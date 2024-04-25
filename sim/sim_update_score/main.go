package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/zjl410091917/game-rank/interal/utils"

	"github.com/go-resty/resty/v2"
)

func main() {
	for i := 0; i < 200; i++ {
		sendUpdateScore(&usParams{
			Pid:    uint64(i + 1),
			Score:  uint32(utils.Between(1, 10000)),
			Level:  10,
			Name:   fmt.Sprintf("p%d", i+1),
			Active: "gr",
		})
	}
}

type usParams struct {
	Pid    uint64 `json:"pid"`
	Score  uint32 `json:"score"`
	Level  uint32 `json:"level"`
	Name   string `json:"name"`
	Active string `json:"active"`
}

func sendUpdateScore(params *usParams) {
	body, _ := json.Marshal(params)
	apiUrl := new(url.URL)
	apiUrl.Scheme = "http"
	apiUrl.Host = "127.0.0.1:6818"
	apiUrl.Path = "api/update_score"
	_, _ = resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(apiUrl.String())
}
