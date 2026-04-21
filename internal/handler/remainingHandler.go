package handler

import (
	"net/http"

	"shortener/internal/middleware"
	"shortener/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// RemainingHandler GET /api/convert/remaining
func RemainingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, userID, _ := middleware.GetUserFromContext(r.Context())

		if userID > 0 {
			httpx.OkJson(w, map[string]int{"remaining": -1}) // -1 means unlimited
			return
		}

		count, err := svcCtx.ShortUrlModel.CountUnregistered(r.Context())
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		remaining := 3 - int(count)
		if remaining < 0 {
			remaining = 0
		}
		httpx.OkJson(w, map[string]int{"remaining": remaining})
	}
}
