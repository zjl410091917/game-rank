package httpx

import (
	"fmt"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/zjl410091917/game-rank/interal/config"
	"github.com/zjl410091917/game-rank/interal/httpx/fiber_cors"
	"github.com/zjl410091917/game-rank/interal/httpx/http_header"
	"github.com/zjl410091917/game-rank/interal/logger"
	"github.com/zjl410091917/game-rank/interal/utils"
	"go.uber.org/zap"
)

type (
	fiberServer struct {
		handlers map[string]HttpHandler
		app      *fiber.App
		tn       atomic.Int64 // task 状态计数
		state    atomic.Bool  // 是否可以访问
	}

	fiberContext struct {
		ctx *fiber.Ctx
	}
)

func newFiberServer() HttpServer {
	return &fiberServer{handlers: make(map[string]HttpHandler)}
}

func (fs *fiberServer) AddHandler(path string, handler HttpHandler) {
	_, ok := fs.handlers[path]
	if ok {
		logger.ErrorWithPanic(fmt.Sprintf("%s handler repeat", path))
	}
	fs.handlers[path] = handler
}

// DoTask 请求任务是否可以执行，可以的话进行计数
func (fs *fiberServer) DoTask() (can bool) {
	can = fs.state.Load()
	if can {
		fs.tn.Add(1)
	}
	return
}

func (fs *fiberServer) TaskFinish() {
	fs.tn.Add(-1)
}

func (fs *fiberServer) OnInit() {
	logger.Info("init", zap.String("module", "FiberServer"))
	fs.app = fiber.New(fiber.Config{
		IdleTimeout:  time.Minute * 5,
		ReadTimeout:  time.Minute * 2,
		WriteTimeout: time.Minute * 2,
	})
	fs.app.Use(printFiberHeader)
	fs.app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	fs.app.Use(recover.New())
	fs.app.Use(fiber_cors.New(fiber_cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Request-Type,Client-Version,Config-Version,Resource-Version",
		AllowCredentials: true,
	}))
	fs.app.Name("Game Rank")

	for key, fun := range fs.handlers {
		fs.app.All(key, toFiberHandler(fun))
	}
}

func toFiberHandler(fun HttpHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fun(newFiberContext(c))
	}
}

func (fs *fiberServer) OnStart() {
	logger.Info("OnStart", zap.String("module", "FiberServer"))
	if fs.state.CompareAndSwap(false, true) {
		err := fs.app.Listen(fmt.Sprintf(":%d", config.C().Server.Port))
		utils.IfErrPanic(err)
	}
}

// OnStop
// 1. k8s环境中，在服务更新替换时，当执行关掉旧服务时，理论上不再回有新请求到达该节点
// 这里依旧停滞2s 最大限度保证切换平滑
// 2. 启动状态置位，后不再接受新的请求。
// 3. 对于已经在处理中的请求，预留一段时间完成，超时后强制关闭
// 4. 服务被执行关闭时，首先被关闭的是httpx服务，只有它退出后，其他服务才能进入退出流程
func (fs *fiberServer) OnStop() {
	time.Sleep(time.Second * 2)
	fs.state.CompareAndSwap(true, false)

	n := 0
	for fs.tn.Load() > 0 && n < 60 {
		time.Sleep(time.Millisecond * 50)
		n++
	}
	logger.Info("OnStop", zap.String("module", "FiberServer"), zap.Int("n", n))
	utils.IfErrPanic(fs.app.Shutdown())
}

func newFiberContext(ctx *fiber.Ctx) HttpContext {
	return &fiberContext{ctx: ctx}
}

func (fctx *fiberContext) GetBody() []byte {
	return fctx.ctx.Body()
}

func (fctx *fiberContext) GetFormValues() url.Values {
	if !http_header.Match(fctx.ctx.GetReqHeaders(), "Content-Type", "application/x-www-form-urlencoded") {
		return map[string][]string{}
	}
	values, err := url.ParseQuery(string(fctx.ctx.Body()))
	if err != nil {
		return map[string][]string{}
	}

	return values
}

func (fctx *fiberContext) GetReqHeaders() map[string][]string {
	return fctx.ctx.GetReqHeaders()
}

func (fctx *fiberContext) Send(data []byte) error {
	return fctx.ctx.Send(data)
}

func (fctx *fiberContext) SendStatus(status int) error {
	return fctx.ctx.SendStatus(status)
}

func printFiberHeader(c *fiber.Ctx) (err error) {
	logger.Info("http-request", zap.Any("header", c.GetReqHeaders()), zap.Any("Port", c.Port()))
	err = c.Next()
	if err != nil {
		return
	}
	logger.Info("http-response", zap.Any("header", func() map[string]any {
		rh := make(map[string]any)
		c.Response().Header.VisitAll(func(key, value []byte) {
			rh[string(key)] = string(value)
		})
		return rh
	}()))
	return
}
