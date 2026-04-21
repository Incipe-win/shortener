package handler

import (
	"net/http"

	"shortener/internal/middleware"
	"shortener/internal/svc"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	c := serverCtx.Config

	cors := middleware.CORSMiddleware(c.CORS.AllowOrigins)
	jwtAuth := middleware.JWTCookieMiddleware(c.Auth.JWTSecret)

	// 注册健康检查端点
	server.AddRoutes([]rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/health",
			Handler: HealthHandler(),
		},
	})

	// 注册限流中间件（基于 Redis，每秒最多 20 次请求）
	rateLimit := cors
	if serverCtx.Config.CacheRedis[0].Pass != "" || serverCtx.Config.CacheRedis[0].Host != "" {
		redisStore := serverCtx.Config.CacheRedis[0]
		r := redis.New(redisStore.Host, func(r *redis.Redis) {
			r.Type = redis.NodeType
			r.Pass = redisStore.Pass
		})
		rateLimit = middleware.RateLimitMiddleware(r, 1, 20)
	}

	// 公开路由（带 CORS + 限流）
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/preview/:short_url",
				Handler: rateLimit(PreviewHandler(serverCtx)),
			},
			{
				Method:  http.MethodPost,
				Path:    "/convert",
				Handler: rateLimit(ConvertHandler(serverCtx)),
			},
			{
				Method:  http.MethodGet,
				Path:    "/:short_url",
				Handler: ShowHandler(serverCtx),
			},
		},
	)

	// Auth 路由（公开 + CORS）
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/auth/login",
				Handler: cors(LoginHandler(c)),
			},
			{
				Method:  http.MethodPost,
				Path:    "/auth/logout",
				Handler: cors(LogoutHandler()),
			},
			{
				Method:  http.MethodOptions,
				Path:    "/auth/login",
				Handler: cors(func(w http.ResponseWriter, r *http.Request) {}),
			},
			{
				Method:  http.MethodOptions,
				Path:    "/auth/logout",
				Handler: cors(func(w http.ResponseWriter, r *http.Request) {}),
			},
			{
				Method:  http.MethodOptions,
				Path:    "/auth/me",
				Handler: cors(func(w http.ResponseWriter, r *http.Request) {}),
			},
			{
				Method:  http.MethodOptions,
				Path:    "/links",
				Handler: cors(func(w http.ResponseWriter, r *http.Request) {}),
			},
		},
	)

	// 需要认证的路由（CORS + JWT）
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/auth/me",
				Handler: cors(jwtAuth(MeHandler())),
			},
			{
				Method:  http.MethodGet,
				Path:    "/links",
				Handler: cors(jwtAuth(LinksHandler(serverCtx))),
			},
		},
	)
}

// HealthHandler GET /health
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpx.OkJson(w, map[string]string{"status": "ok"})
	}
}
