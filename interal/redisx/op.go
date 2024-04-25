package redisx

import (
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
)

// R 获取默认连接
func R() redis.UniversalClient {
	return GetInstance().conList[0].client
}

func ID(pid uint64) string {
	index := pid % uint64(len(instance.conList))
	return instance.conList[index].id
}

// PR 通过id 获取redis链接
func PR(pid uint64) redis.UniversalClient {
	index := pid % uint64(len(instance.conList))
	return instance.conList[index].client
}

func C(id string) redis.UniversalClient {
	return instance.conMap[id].client
}

func NewLock(mutexName string) *redsync.Mutex {
	rs := GetInstance().rs
	return rs.NewMutex(mutexName, redsync.WithExpiry(time.Second*10))
}
