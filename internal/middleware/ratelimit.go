package middleware

import (
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// RateLimitMiddleware 基于 IP 的请求限流
// period: 窗口时长（秒），quota: 窗口内最大请求数
func RateLimitMiddleware(store *redis.Redis, period int, quota int) func(http.HandlerFunc) http.HandlerFunc {
	limiter := limit.NewPeriodLimit(period, quota, store, "shortener_rate_limit")

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ip := extractIP(r)
			key := "ip:" + ip

			code, err := limiter.Take(key)
			if err != nil {
				next(w, r) // Redis 不可用时放行
				return
			}

			if code == limit.OverQuota {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"code":429,"msg":"请求过于频繁，请稍后重试"}`))
				return
			}

			next(w, r)
		}
	}
}

// extractIP 从请求中提取真实 IP
func extractIP(r *http.Request) string {
	// Nginx 反向代理设置
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// 取第一个 IP
		if idx := strings.Index(ip, ","); idx != -1 {
			return strings.TrimSpace(ip[:idx])
		}
		return ip
	}
	return r.RemoteAddr
}
