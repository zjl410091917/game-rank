package entry

type ActiveRank struct {
	ActiveId  string        // 活动id
	MaxNumber int32         // 最大编号
	Buckets   []*RankBucket // 排名桶信息
}

type RankBucket struct {
	Number  int32     // 序号，与bucket之间的排序无关，用于拼mongo集合名
	Count   int32     // 当前拥有玩家数据量
	MinRank *RankInfo // 最小排名信息
	MaxRank *RankInfo // 最大排名信息
}

func (a *ActiveRank) Bucket(bi int) (b *RankBucket) {
	switch {
	case len(a.Buckets) > bi:
		b = a.Buckets[bi]
	case len(a.Buckets) == bi:
		a.MaxNumber++
		a.Buckets = append(a.Buckets, &RankBucket{
			Number:  a.MaxNumber,
			Count:   0,
			MinRank: nil,
			MaxRank: nil,
		})
		b = a.Buckets[bi]
	}
	return
}
