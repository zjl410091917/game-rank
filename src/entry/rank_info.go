package entry

import "fmt"

type RankInfo struct {
	Level        uint32 // 等级
	Score        uint32 // 积分
	Rank         uint32 // 全局排名
	BRank        uint32 // 桶排名
	Pid          uint64 // 用户id
	UpdateTime   int64  // 更新时间
	Name         string // 姓名
	BucketNumber int32  // 桶编号
	BucketIndex  int    // 桶索引
	GroupId      uint64 // 排名group id
}

type RankSed struct {
	Pid uint64 // 用户id
	Sed string // 排名系数: "`100000 + score` + time + level + name" 在mongo内进行字符串排序
}

// CRankInfo 前端返回值单个信息
type CRankInfo struct {
	Level uint32 `json:"level,omitempty"` // 等级
	Score uint32 `json:"score,omitempty"` // 积分
	Rank  uint32 `json:"rank,omitempty"`  // 全局排名
	Pid   uint64 `json:"pid,omitempty"`   // 用户id
	Name  string `json:"name,omitempty"`  // 姓名
}

func (r *RankInfo) BI() int {
	return max(r.BucketIndex-1, 0)
}

func NewRankSed(ri *RankInfo) *RankSed {
	rs := &RankSed{
		Pid: ri.Pid,
		Sed: Sed(ri),
	}
	return rs
}

func Sed(ri *RankInfo) string {
	return fmt.Sprintf("%d%d%d%s", ri.Score+100000, ri.UpdateTime, ri.Level+10000, ri.Name)
}
