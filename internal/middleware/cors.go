package middleware

import "net/http"

// CORSMiddleware 跨域中间件，支持 HTTP-only Cookie 认证
func CORSMiddleware(allowOrigins []string) func(http.HandlerFunc) http.HandlerFunc {
	originSet := make(map[string]struct{}, len(allowOrigins))
	for _, o := range allowOrigins {
		originSet[o] = struct{}{}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// 检查是否在允许列表中
			if _, ok := originSet[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, traceparent, tracestate")
				w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			// Preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next(w, r)
		}
	}
}
