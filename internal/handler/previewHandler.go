package handler

import (
	"net/http"
	"time"

	"shortener/internal/logic"
	"shortener/internal/svc"
	"shortener/internal/types"
	"shortener/pkg/metrics"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PreviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			metrics.RequestDuration.WithLabelValues("GET", "/preview/:short_url").Observe(time.Since(start).Seconds())
		}()

		var req types.PreviewRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewPreviewLogic(r.Context(), svcCtx)
		resp, err := l.Preview(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
