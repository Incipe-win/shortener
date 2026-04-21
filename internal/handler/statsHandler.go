package handler

import (
	"net/http"

	"shortener/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func StatsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := ctx.ShortUrlModel.GetStats(r.Context())
		if err != nil {
			logx.Errorw("failed to get stats", logx.LogField{Key: "err", Value: err.Error()})
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, stats)
	}
}
