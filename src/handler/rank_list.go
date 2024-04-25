package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zjl410091917/game-rank/interal/redisx"

	"github.com/zjl410091917/game-rank/src/logic"

	"github.com/zjl410091917/game-rank/src/dao"
	"github.com/zjl410091917/game-rank/src/entry"
)

type rlParams struct {
	Pid    uint64 `json:"pid"`
	Active string `json:"active"`
}

type rlRes struct {
	RankList []*entry.CRankInfo `json:"rankList,omitempty"`
	RIndex   int                `json:"rIndex,omitempty"`
}

func onRequestRankList(params []byte) (res *rlRes, err error) {
	param := &rlParams{}
	err = json.Unmarshal(params, param)
	if err != nil {
		return
	}

	res = &rlRes{}
	rk := fmt.Sprintf("rank_list_%d", param.Pid)
	ok, err := dao.GetCache(context.Background(), redisx.PR(param.Pid), rk, res)
	if err != nil || ok {
		return
	}
	ar := &entry.ActiveRank{
		ActiveId:  param.Active,
		MaxNumber: 0,
		Buckets:   nil,
	}
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
	rl, pi, err := logic.GetRankList(ar, ri)
	if err != nil {
		return
	}

	res.RankList = rl
	res.RIndex = pi
	err = dao.SetCache(context.Background(), redisx.PR(param.Pid), rk, res)
	return
}
