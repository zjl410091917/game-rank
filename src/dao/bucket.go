package dao

import (
	"context"
	"fmt"

	"github.com/zjl410091917/game-rank/src/common"

	"github.com/zjl410091917/game-rank/interal/mongox"
	"github.com/zjl410091917/game-rank/interal/redisx"
	"github.com/zjl410091917/game-rank/src/entry"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetActiveRank(ctx context.Context, ar *entry.ActiveRank) (err error) {
	rk := fmt.Sprintf("active_rank_%s", ar.ActiveId)
	ok, err := GetCache(ctx, redisx.R(), rk, ar)
	if err != nil || ok {
		return
	}
	err = mongox.FindOne(ctx, "ActiveRank", bson.M{"aid": ar.ActiveId}, ar)
	if err != nil {
		return
	}

	err = SetCache(ctx, redisx.R(), rk, ar)
	return
}

func UpdateActiveRank(ctx context.Context, ar *entry.ActiveRank) (err error) {
	err = mongox.UpdateOne(ctx, "ActiveRank", bson.M{"aid": ar.ActiveId}, ar, true)
	if err != nil {
		return
	}
	err = SetCache(ctx, redisx.R(), fmt.Sprintf("active_rank_%s", ar.ActiveId), ar)
	return
}

func UpdateBucketSed(ctx context.Context, bn int32, rs *entry.RankSed) error {
	return mongox.UpdateOne(ctx, fmt.Sprintf("Bucket-%d", bn), bson.M{"pid": rs.Pid}, rs, true)
}

func DeleteBucketSed(ctx context.Context, bn int32, pid uint64) (err error) {
	_, err = mongox.DeleteOne(ctx, fmt.Sprintf("Bucket-%d", bn), bson.M{"pid": pid})
	return
}

// FindMinBucketSed 获取bucket最小值
func FindMinBucketSed(ctx context.Context, bn int32) (sed *entry.RankSed, err error) {
	sed = &entry.RankSed{}
	opts := options.FindOne()
	opts.SetSort(bson.M{"sed": 1})
	err = mongox.FindOneWithOpts(ctx, fmt.Sprintf("Bucket-%d", bn), bson.M{}, sed, opts)
	return
}

func FindTopBucketSed(ctx context.Context, bn int32, limit int64) (sedList []*entry.RankSed, err error) {
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSort(bson.M{"sed": -1})
	err = mongox.Find(ctx, fmt.Sprintf("Bucket-%d", bn), bson.M{}, &sedList, opts)
	return
}

func DeleteSomeBucketSed(ctx context.Context, bn int32, pidList []uint64) (err error) {
	_, err = mongox.DeleteMany(ctx, fmt.Sprintf("Bucket-%d", bn), bson.M{"pid": bson.M{"$in": pidList}})
	return
}

func InsertBucketSed(ctx context.Context, bn int32, sedList []*entry.RankSed) (err error) {
	var list []any
	for i := 0; i < len(sedList); i++ {
		list = append(list, sedList[i])
	}
	err = mongox.Insert(ctx, fmt.Sprintf("Bucket-%d", bn), list)
	return
}

func CalRankInBucket(ctx context.Context, sed string, bn int32) (n int64, err error) {
	n, err = mongox.Count(ctx, fmt.Sprintf("Bucket-%d", bn), bson.M{"sed": bson.M{"$gte": sed}})
	return
}

func GetSedListInBucket(ctx context.Context, sed string, bn int32, direction int32) (sedList []*entry.RankSed, err error) {
	opts := options.Find()
	opts.SetLimit(common.CRCount)
	opts.SetSort(bson.M{"sed": -1})
	filterK := "$lt"
	filter := bson.M{}
	if direction == common.DLeft {
		opts.SetSort(bson.M{"sed": 1})
		filterK = "$gt"
	}
	if sed != "" {
		filter = bson.M{"sed": bson.M{filterK: sed}}
	}

	err = mongox.Find(ctx, fmt.Sprintf("Bucket-%d", bn), filter, &sedList, opts)
	return
}
