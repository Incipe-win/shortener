package handler

import (
	"net/http"

	"shortener/internal/logic"
	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ConvertHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConvertRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if err := validator.New().StructCtx(r.Context(), &req); err != nil {
			logx.Errorw("validation error", logx.LogField{Key: "error", Value: err.Error()})
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewConvertLogic(r.Context(), svcCtx)
		resp, err := l.Convert(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
