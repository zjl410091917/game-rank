package redisx

import (
	"sync"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"github.com/zjl410091917/game-rank/interal/config"
	"github.com/zjl410091917/game-rank/interal/logger"
	"github.com/zjl410091917/game-rank/interal/utils"
	"go.uber.org/zap"
)

type Conn struct {
	id     string
	client redis.UniversalClient
}

type Redisx struct {
	conList []*Conn
	conMap  map[string]*Conn
	rs      *redsync.Redsync
}

var (
	instance *Redisx
	once     sync.Once
)

func GetInstance() *Redisx {
	once.Do(func() {
		instance = &Redisx{
			conMap: make(map[string]*Conn),
		}
	})
	return instance
}

func (r *Redisx) OnInit() {
	logger.Info("init", zap.String("module", "Redis"))
	cc := config.C()

	for i := 0; i < len(cc.Redis); i++ {
		if !cc.Redis[i].Cluster {
			client := redis.NewClient(&redis.Options{
				Addr: cc.Redis[i].Addr,
				DB:   cc.Redis[i].DB,
			})
			id := utils.MD5(cc.Redis[i].Addr)
			con := &Conn{id: id, client: client}
			r.conList = append(r.conList, con)
			r.conMap[id] = con
		} else {
			client := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs: []string{cc.Redis[i].Addr},
			})
			id := utils.MD5(cc.Redis[i].Addr)
			con := &Conn{id: id, client: client}
			r.conList = append(r.conList, con)
			r.conMap[id] = con
		}
	}
	pool := goredis.NewPool(r.conList[0].client)
	r.rs = redsync.New(pool)
}

func (r *Redisx) OnStart() {
	logger.Info("OnStart", zap.String("module", "Redis"))
}

func (r *Redisx) OnStop() {
	logger.Info("OnStop", zap.String("module", "Redis"))
	for i := 0; i < len(r.conList); i++ {
		_ = r.conList[i].client.Close()
	}
}
