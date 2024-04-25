package httpx

import (
	"sync"
)

var (
	instance HttpServer
	once     sync.Once
)

func GetInstance() HttpServer {
	once.Do(func() {
		instance = newFiberServer()
	})
	return instance
}
