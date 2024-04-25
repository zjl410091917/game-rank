package dao

import (
	"context"
	"fmt"

	"github.com/zjl410091917/game-rank/interal/mongox"
	"github.com/zjl410091917/game-rank/interal/redisx"
	"github.com/zjl410091917/game-rank/src/entry"
	"go.mongodb.org/mongo-driver/bson"
)

func GetRankInfo(ctx context.Context, info *entry.RankInfo) (err error) {
	rk := fmt.Sprintf("rank_info_%d", info.Pid)
	ok, err := GetCache(ctx, redisx.PR(info.Pid), rk, info)
	if err != nil || ok {
		return
	}
	err = mongox.FindOne(ctx, "RankInfo", bson.M{"pid": info.Pid}, info)
	if err != nil {
		return
	}

	err = SetCache(ctx, redisx.PR(info.Pid), rk, info)
	return
}

func UpdateRankInfo(ctx context.Context, info *entry.RankInfo) (err error) {
	err = mongox.UpdateOne(ctx, "RankInfo", bson.M{"pid": info.Pid}, info, true)
	if err != nil {
		return
	}
	err = SetCache(ctx, redisx.PR(info.Pid), fmt.Sprintf("rank_info_%d", info.Pid), info)
	return
}

func ReplaceBNSomeRankInfo(ctx context.Context, pidList []uint64, bn, bi int) (err error) {
	err = mongox.UpdateMany(ctx, "RankInfo", bson.M{"pid": bson.M{"$in": pidList}}, bson.M{"bucketnumber": bn, "bucketindex": bi}, false)
	if err != nil {
		return
	}
	// 删除用户rank的cache缓存
	pk := map[string][]string{}
	for i := 0; i < len(pidList); i++ {
		ri := redisx.ID(pidList[i])
		rk := fmt.Sprintf("rank_info_%d", pidList[i])
		pk[ri] = append(pk[ri], rk)
	}

	for ri, rks := range pk {
		rd := redisx.C(ri)
		err = rd.Del(ctx, rks...).Err()
		if err != nil {
			return
		}
	}
	return
}

func LoadRankInfoList(ctx context.Context, pidList []uint64) (rm map[uint64]*entry.RankInfo, err error) {
	var rl []*entry.RankInfo
	err = mongox.Find(ctx, "RankInfo", bson.M{"pid": bson.M{"$in": pidList}}, &rl)
	if err != nil {
		return
	}
	rm = make(map[uint64]*entry.RankInfo)
	for i := 0; i < len(rl); i++ {
		rm[rl[i].Pid] = rl[i]
	}
	return
}
