package logic

import (
	"context"
	"math"
	"slices"

	"github.com/zjl410091917/game-rank/src/dao"
	"github.com/zjl410091917/game-rank/src/entry"
)

// onBucketAdd 桶内新增一条记录
func onBucketAdd(ar *entry.ActiveRank, info *entry.RankInfo) {
	b := ar.Bucket(info.BI())
	info.BucketNumber = b.Number
	b.Count++
	// 比最大值大
	if compRank(info, b.MaxRank) == crBig {
		b.MaxRank = info
	}

	// 比最小值小
	if compRank(info, b.MinRank) == crSmall {
		b.MinRank = info
	}

	if b.Count == 1 {
		b.MinRank = b.MaxRank
	}
}

// onBucketSedUpdate 桶内一条记录更新了(肯定是积分增加了)
func onBucketSedUpdate(ar *entry.ActiveRank, info *entry.RankInfo) (err error) {
	b := ar.Bucket(info.BI())
	// 比最大的大
	if compRank(info, b.MaxRank) == crBig {
		b.MaxRank = info
	}

	// 重新计算最小值
	if info.Pid == b.MinRank.Pid {
		err = resetBucketMin(b)
	}
	return
}

// resetBucketMin 重置最小值
func resetBucketMin(b *entry.RankBucket) (err error) {
	if b.Count == 1 {
		b.MinRank = b.MaxRank
		return
	}
	sed, err := dao.FindMinBucketSed(context.Background(), b.Number)
	if err != nil {
		return
	}
	b.MinRank.Pid = sed.Pid
	err = dao.GetRankInfo(context.Background(), b.MinRank)
	return
}

// getBucket 获取桶信息
// rank_info中记录的索引可能失效，通过number进行查找，并更新索引值
func getBucket(ar *entry.ActiveRank, info *entry.RankInfo) (b *entry.RankBucket) {
	b = ar.Bucket(info.BI())
	if b.Number == info.BucketNumber {
		return
	}
	for i := 0; i < len(ar.Buckets); i++ {
		if ar.Buckets[i].Number == info.BucketNumber {
			b = ar.Buckets[i]
			info.BucketIndex = i + 1
			break
		}
	}
	return
}

func updateBucketSed(info *entry.RankInfo) (err error) {
	rs := entry.NewRankSed(info)
	err = dao.UpdateBucketSed(context.Background(), info.BucketNumber, rs)
	return
}

func deleteBucketSed(pid uint64, bn int32) (err error) {
	err = dao.DeleteBucketSed(context.Background(), bn, pid)
	return
}

// checkBucketCap 检查桶是否需要分裂
func checkBucketCap(ar *entry.ActiveRank, bi int) (err error) {
	b := ar.Bucket(bi)
	if b.Count < bucketMax {
		return
	}
	err = splitBucket(ar, b, bi)
	return
}

// splitBucket 拆分桶：将top 30%的记录放入新桶
func splitBucket(ar *entry.ActiveRank, b *entry.RankBucket, bi int) (err error) {
	limit := math.Round(float64(b.Count) * 0.3)
	sedList, err := dao.FindTopBucketSed(context.Background(), b.Number, int64(limit))
	if err != nil {
		return
	}
	err = dao.InsertBucketSed(context.Background(), ar.MaxNumber+1, sedList[0:len(sedList)-1])
	if err != nil {
		return
	}
	pidList := make([]uint64, 0, len(sedList))
	for i := 0; i < len(sedList)-1; i++ {
		pidList = append(pidList, sedList[i].Pid)
	}
	err = dao.DeleteSomeBucketSed(context.Background(), b.Number, pidList)
	if err != nil {
		return
	}
	// 更新rank_info
	err = dao.ReplaceBNSomeRankInfo(context.Background(), pidList, int(ar.MaxNumber+1), bi)
	if err != nil {
		return
	}
	b.Count -= int32(len(pidList))
	b.MaxRank.Pid = sedList[len(sedList)-1].Pid
	err = dao.GetRankInfo(context.Background(), b.MaxRank)
	if err != nil {
		return
	}
	nb := &entry.RankBucket{
		Number:  ar.MaxNumber + 1,
		Count:   int32(len(pidList)),
		MinRank: &entry.RankInfo{},
		MaxRank: &entry.RankInfo{},
	}
	nb.MaxRank.Pid = pidList[0]
	nb.MinRank.Pid = pidList[len(pidList)-1]
	err = dao.GetRankInfo(context.Background(), nb.MinRank)
	if err != nil {
		return
	}
	err = dao.GetRankInfo(context.Background(), nb.MaxRank)
	if err != nil {
		return
	}
	ar.MaxNumber++
	ar.Buckets = slices.Insert(ar.Buckets, bi, nb)
	return
}
