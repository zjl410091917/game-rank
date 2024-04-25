package app

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/zjl410091917/game-rank/interal/config"

	"go.uber.org/zap"
)

type (
	Module interface {
		OnInit()
		OnStart()
		OnStop()
	}
	App struct {
		status     int32
		shutdown   chan int
		moduleList []Module
		name       string
		logger     *zap.Logger
	}
)

var (
	_app *App
	once sync.Once
)

func GetInstance() *App {
	once.Do(func() {
		_app = &App{
			shutdown: make(chan int),
		}
	})
	return _app
}

func (app *App) SetName(n string) *App {
	app.name = n
	config.GetInstance().Name = n
	return app
}

func (app *App) Name() string {
	return app.name
}

func (app *App) SetLogger(l *zap.Logger) *App {
	app.logger = l
	return app
}

func (app *App) Run() {
	if !atomic.CompareAndSwapInt32(&app.status, 0, 1) {
		panic(fmt.Errorf("app run repeat"))
	}
	go app.processSignal()

	app.onRun()
	if app.logger != nil {
		app.logger.Info("run",
			zap.String("Name", app.Name()),
			zap.String("Goroutine", strconv.Itoa(runtime.NumGoroutine())),
			zap.String("CPU", strconv.Itoa(runtime.NumCPU())),
		)
	}
	// bp, _ := os.Getwd()
	// dtalk.SendMessage(fmt.Sprintf("* **Name**: %s\n\n* **Branch**: %s\n\n* **Commit**: %s\n\n* **Deploy**: %s\n\n* **Env**: %s", app.name, version.Branch, version.Commit, bp, config.GetInstance().Env), "Launch Message")
	<-app.shutdown
	app.onShutdown()
}

func (app *App) Module(m Module) *App {
	if atomic.CompareAndSwapInt32(&app.status, 1, 0) {
		panic(fmt.Errorf("app is running"))
	}
	app.moduleList = append(app.moduleList, m)
	m.OnInit()
	return app
}

func (app *App) processSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown <- 1
}

func (app *App) onRun() {
	for i := 0; i < len(app.moduleList); i++ {
		go app.moduleList[i].OnStart()
	}
}

func (app *App) onShutdown() {
	// time.Sleep(time.Second * 3) // 为当前正在处理的任务留时间
	if atomic.CompareAndSwapInt32(&app.status, 1, 0) {
		for i := len(app.moduleList) - 1; i >= 0; i-- {
			app.moduleList[i].OnStop()
		}
	}
}
