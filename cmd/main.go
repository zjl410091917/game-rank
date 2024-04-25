package main

import (
	"github.com/zjl410091917/game-rank/interal/app"
	"github.com/zjl410091917/game-rank/interal/config"
	"github.com/zjl410091917/game-rank/interal/httpx"
	"github.com/zjl410091917/game-rank/interal/logger"
	"github.com/zjl410091917/game-rank/interal/mongox"
	"github.com/zjl410091917/game-rank/interal/redisx"
	_ "github.com/zjl410091917/game-rank/src/handler"
	_ "go.uber.org/automaxprocs"
)

func main() {
	a := app.GetInstance()
	a.SetName("server").
		Module(config.GetInstance()).
		Module(logger.GetInstance()).
		Module(mongox.GetInstance()).
		Module(redisx.GetInstance()).
		Module(httpx.GetInstance()).
		Run()
}
