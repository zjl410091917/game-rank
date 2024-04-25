package mongox

import (
	"context"
	"fmt"
	"sync"

	"github.com/zjl410091917/game-rank/interal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongox struct {
	conn *mongo.Client
}

var (
	instance *Mongox
	once     sync.Once
)

func GetInstance() *Mongox {
	once.Do(func() {
		instance = &Mongox{}
	})
	return instance
}

func (m *Mongox) OnInit() {
}

func (m *Mongox) OnStart() {
	cc := config.C().Mongo
	uri := fmt.Sprintf("mongodb://%s:%d/%s?w=majority", cc.Host, cc.Port, cc.DB)
	o := options.Client().ApplyURI(uri)
	o.SetMaxPoolSize(100)
	mc, err := mongo.Connect(context.Background(), o)
	if err != nil {
		panic(fmt.Errorf("connect %s error", uri))
	}
	instance.conn = mc
}

func (m *Mongox) OnStop() {
	_ = m.conn.Disconnect(context.Background())
}

func (m *Mongox) collection(cname string) *mongo.Collection {
	return m.conn.Database(config.C().Mongo.DB).Collection(cname)
}
