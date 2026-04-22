package handler

import (
	"context"
	"net/http"
	"time"

	"shortener/internal/ctxdata"
	"shortener/internal/logic"
	"shortener/internal/svc"
	"shortener/internal/types"
	"shortener/pkg/metrics"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ShowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			metrics.RequestDuration.WithLabelValues("GET", "/:short_url").Observe(time.Since(start).Seconds())
		}()

		var req types.ShowRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 将 HTTP 请求信息注入 context，供 ShowLogic 发送点击事件
		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxdata.KeyClientIP, clientIP(r))
		ctx = context.WithValue(ctx, ctxdata.KeyUserAgent, r.UserAgent())
		ctx = context.WithValue(ctx, ctxdata.KeyReferer, r.Referer())

		l := logic.NewShowLogic(ctx, svcCtx)
		resp, err := l.Show(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			// 如果有安全警告，添加到响应头中
			if resp.RiskWarning != "" {
				w.Header().Set("X-Safety-Warning", resp.RiskWarning)
			}
			http.Redirect(w, r, resp.LongUrl, http.StatusFound)
		}
	}
}

// clientIP 提取客户端真实 IP
func clientIP(r *http.Request) string {
	// 优先从常见代理头获取
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
