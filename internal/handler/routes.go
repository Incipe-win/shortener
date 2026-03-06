package handler

import (
	"net/http"

	"shortener/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/preview/:short_url",
				Handler: PreviewHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/convert",
				Handler: ConvertHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/:short_url",
				Handler: ShowHandler(serverCtx),
			},
		},
	)
}
