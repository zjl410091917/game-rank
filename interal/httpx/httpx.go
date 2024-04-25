package httpx

import (
	"net/url"

	"github.com/zjl410091917/game-rank/interal/app"
)

type (
	HttpHandler    func(HttpContext) error
	HttpErrHandler func(HttpContext, string) error

	HttpServer interface {
		app.Module
		AddHandler(path string, handler HttpHandler)
		DoTask() bool
		TaskFinish()
	}

	HttpContext interface {
		GetBody() []byte
		GetFormValues() url.Values
		GetReqHeaders() map[string][]string
		SendStatus(status int) error
		Send(data []byte) error
	}
)
