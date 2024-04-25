package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zjl410091917/game-rank/interal/redisx"

	"github.com/zjl410091917/game-rank/src/logic"

	"github.com/zjl410091917/game-rank/src/dao"
	"github.com/zjl410091917/game-rank/src/entry"
)

type updateParam struct {
	Pid    uint64 `json:"pid"`
	Score  uint32 `json:"score"`
	Level  uint32 `json:"level"`
	Name   string `json:"name"`
	Active string `json:"active"`
}

func onRequestUpdateScore(params []byte) (err error) {
	param := &updateParam{}
	err = json.Unmarshal(params, param)
	if err != nil {
		return
	}

	ar := &entry.ActiveRank{
		ActiveId:  param.Active,
		MaxNumber: 0,
		Buckets:   nil,
	}
	rl := redisx.NewLock("update_score")
	err = rl.Lock()
	if err != nil {
		return
	}
	defer func() {
		_, err = rl.Unlock()
	}()

	ctx := context.Background()
	err = dao.GetActiveRank(ctx, ar)
	if err != nil {
		return
	}
	ri := &entry.RankInfo{
		Pid: param.Pid,
	}
	err = dao.GetRankInfo(ctx, ri)
	if err != nil {
		return
	}
	// 低于当前分数
	if param.Score < ri.Score {
		return
	}
	ri.Level = param.Level
	ri.Name = param.Name
	ri.UpdateTime = time.Now().UnixMilli()
	ri.Score = param.Score
	err = logic.UpdateScore(ar, ri)
	return
}
