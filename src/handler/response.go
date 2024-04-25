package handler

import (
	"encoding/json"

	"github.com/zjl410091917/game-rank/interal/logger"
	"go.uber.org/zap"

	"github.com/zjl410091917/game-rank/interal/httpx"
)

type response struct {
	Status int    `json:"status"`
	Data   any    `json:"data,omitempty"`
	Err    string `json:"err,omitempty"`
}

func setResponse(ctx httpx.HttpContext, data any, resErr error) error {
	res := &response{
		Status: 1,
		Data:   data,
		Err:    "",
	}
	if resErr != nil {
		res.Err = resErr.Error()
		res.Status = 0
	}
	body, err := json.Marshal(res)
	if err != nil {
		logger.Error("setResponse", zap.Error(err))
		return ctx.SendStatus(500)
	}
	return ctx.Send(body)
}
