package handler

import "github.com/zjl410091917/game-rank/interal/httpx"

func init() {
	httpx.GetInstance().AddHandler("/api/update_score", func(ctx httpx.HttpContext) error {
		return setResponse(ctx, nil, onRequestUpdateScore(ctx.GetBody()))
	})

	httpx.GetInstance().AddHandler("/api/rank_list", func(ctx httpx.HttpContext) error {
		rl, err := onRequestRankList(ctx.GetBody())
		return setResponse(ctx, rl, err)
	})
}
