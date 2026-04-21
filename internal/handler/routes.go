package handler

import (
	"net/http"

	"shortener/internal/middleware"
	"shortener/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	c := serverCtx.Config

	cors := middleware.CORSMiddleware(c.CORS.AllowOrigins)
	jwtAuth := middleware.JWTCookieMiddleware(c.Auth.JWTSecret)

	// 公开路由（带 CORS）
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/preview/:short_url",
				Handler: cors(PreviewHandler(serverCtx)),
			},
			{
				Method:  http.MethodPost,
				Path:    "/convert",
				Handler: cors(ConvertHandler(serverCtx)),
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
				Path:    "/api/auth/login",
				Handler: cors(LoginHandler(c)),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/auth/logout",
				Handler: cors(LogoutHandler()),
			},
			{
				Method:  http.MethodOptions,
				Path:    "/api/auth/login",
				Handler: cors(func(w http.ResponseWriter, r *http.Request) {}),
			},
			{
				Method:  http.MethodOptions,
				Path:    "/api/auth/logout",
				Handler: cors(func(w http.ResponseWriter, r *http.Request) {}),
			},
			{
				Method:  http.MethodOptions,
				Path:    "/api/auth/me",
				Handler: cors(func(w http.ResponseWriter, r *http.Request) {}),
			},
			{
				Method:  http.MethodOptions,
				Path:    "/api/links",
				Handler: cors(func(w http.ResponseWriter, r *http.Request) {}),
			},
		},
	)

	// 需要认证的路由（CORS + JWT）
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/auth/me",
				Handler: cors(jwtAuth(MeHandler())),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/links",
				Handler: cors(jwtAuth(LinksHandler(serverCtx))),
			},
		},
	)
}
