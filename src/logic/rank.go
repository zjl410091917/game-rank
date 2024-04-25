package logic

import (
	"context"

	"github.com/zjl410091917/game-rank/src/common"

	"github.com/zjl410091917/game-rank/src/dao"
	"github.com/zjl410091917/game-rank/src/entry"
)

// UpdateScore 更新积分
// info中已经更新了score等相关值，只有bucketIndex待确定
func UpdateScore(ar *entry.ActiveRank, info *entry.RankInfo) (err error) {
	oldNum := info.BucketNumber
	newRecord := info.BucketIndex == 0
	if newRecord {
		onNewScoreInsert(ar, info)
	} else {
		onOldScoreUpdate(ar, info)
	}
	newNum := info.BucketNumber

	// 更新rank_info
	err = dao.UpdateRankInfo(context.Background(), info)
	if err != nil {
		return
	}

	// 根据newNum 更新bucket_sed
	err = updateBucketSed(info)
	if err != nil {
		return
	}

	err = onBucketSedUpdate(ar, info)
	if err != nil {
		return
	}

	if !newRecord && oldNum != newNum {
		// 如果oldNum不等于newNum,即换桶了，则删除之前桶内的记录
		err = deleteBucketSed(info.Pid, oldNum)
		if err != nil {
			return
		}
	}
	err = checkBucketCap(ar, info.BI())
	if err != nil {
		return
	}
	err = dao.UpdateActiveRank(context.Background(), ar)
	return
}

// GetRankList 获取用户排名列表
// 1. 检查是否有最终缓存
// 2. 是否有group缓存
// 3. 重建缓存
func GetRankList(ar *entry.ActiveRank, info *entry.RankInfo) (rl []*entry.CRankInfo, pi int, err error) {
	err = calPlayerRank(ar, info)
	if err != nil {
		return
	}
	err = dao.UpdateRankInfo(context.Background(), info)
	if err != nil {
		return
	}
	rl, pi, err = loadRankList(ar, info)
	return
}

// calPlayerRank 计算玩家的桶内排名和全局排名
func calPlayerRank(ar *entry.ActiveRank, info *entry.RankInfo) (err error) {
	if info.Rank > 0 {
		return
	}
	sed := entry.Sed(info)
	rankInBucket, err := dao.CalRankInBucket(context.Background(), sed, info.BucketNumber)
	if err != nil {
		return
	}
	info.BRank = uint32(rankInBucket)
	info.Rank += uint32(rankInBucket)
	for i := 0; i < len(ar.Buckets); i++ {
		if ar.Buckets[i].Number == info.BucketNumber {
			break
		}
		info.Rank += uint32(ar.Buckets[i].Count)
	}
	return
}

// loadRankList 加载玩家排名附近的信息
func loadRankList(ar *entry.ActiveRank, info *entry.RankInfo) (rl []*entry.CRankInfo, pi int, err error) {
	sed := entry.Sed(info)
	sedRList, err := dao.GetSedListInBucket(context.Background(), sed, info.BucketNumber, common.DRight)
	if err != nil {
		return
	}
	_ = getBucket(ar, info)
	var temp []*entry.RankSed
	if len(sedRList) < common.CRCount && len(ar.Buckets) > info.BucketIndex {
		temp, err = dao.GetSedListInBucket(context.Background(), "", ar.Buckets[info.BucketIndex].Number, common.DRight)
		if err != nil {
			return
		}
		for i := 0; i < len(temp) && len(sedRList) < common.CRCount; i++ {
			sedRList = append(sedRList, temp[i])
		}
	}

	sedLList, err := dao.GetSedListInBucket(context.Background(), sed, info.BucketNumber, common.DLeft)
	if err != nil {
		return
	}
	if len(sedLList) < common.CRCount && info.BucketIndex > 1 {
		temp, err = dao.GetSedListInBucket(context.Background(), "", ar.Buckets[info.BucketIndex-2].Number, common.DLeft)
		if err != nil {
			return
		}
		for i := 0; i < len(temp) && len(sedLList) < common.CRCount; i++ {
			sedLList = append(sedLList, temp[i])
		}
	}

	var pidList []uint64
	for i := len(sedLList) - 1; i >= 0; i-- {
		pidList = append(pidList, sedLList[i].Pid)
	}
	pi = len(pidList)
	pidList = append(pidList, info.Pid)
	for i := 0; i < len(sedRList); i++ {
		pidList = append(pidList, sedRList[i].Pid)
	}

	// 加载rank_info
	rm, err := dao.LoadRankInfoList(context.Background(), pidList)
	if err != nil {
		return
	}
	for i := 0; i < len(pidList); i++ {
		ri := &entry.CRankInfo{
			Level: rm[pidList[i]].Level,
			Score: rm[pidList[i]].Score,
			Rank:  info.Rank,
			Pid:   rm[pidList[i]].Pid,
			Name:  rm[pidList[i]].Name,
		}
		if pi >= i {
			ri.Rank -= uint32(pi - i)
		} else {
			ri.Rank += uint32(i - pi)
		}
		rl = append(rl, ri)
	}
	return
}

// onNewScoreInsert 新积分记录 确定要进入的桶
func onNewScoreInsert(ar *entry.ActiveRank, info *entry.RankInfo) {
	if len(ar.Buckets) == 0 {
		info.BucketIndex = 1 // 默认在第一个桶里
		onBucketAdd(ar, info)
		return
	}
	for i := len(ar.Buckets) - 1; i >= 0; i-- {
		b := ar.Buckets[i]
		info.BucketIndex = i + 1
		if compRank(b.MaxRank, info) == crBig {
			break
		}
	}
	onBucketAdd(ar, info)
}

// onOldScoreUpdate 旧积分记录被更新了(只能是积分增加)
// 1. 与最大值进行比较
// 2. 所在桶是第一个桶，或者不比最大值大，则依然在这个桶内
// 3. 否则寻找更高级的桶，同时删除该桶内的bucket_sed
func onOldScoreUpdate(ar *entry.ActiveRank, info *entry.RankInfo) {
	b := getBucket(ar, info)
	cr := compRank(info, b.MaxRank)
	if info.BucketIndex == 1 || cr != crBig {
		return // 依然在这个桶内
	}
	// 删除桶内的bucket_sed
	b.Count--

	// 寻找更高级的桶
	for i := info.BucketIndex - 2; i >= 0; i-- {
		b = ar.Buckets[i]
		info.BucketIndex = i + 1
		if compRank(b.MaxRank, info) == crBig {
			break
		}
	}
	onBucketAdd(ar, info)
}

// compRank 比较大小，a肯定有值
func compRank(a, b *entry.RankInfo) (res int) {
	if b == nil {
		res = crBig
		return
	}

	aSed := entry.Sed(a)
	bSed := entry.Sed(b)
	switch {
	case aSed > bSed:
		res = crBig
	case aSed < bSed:
		res = crSmall
	default:
		res = crEqual
	}
	return
}
