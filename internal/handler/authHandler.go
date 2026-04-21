package handler

import (
	"net/http"
	"time"

	"shortener/internal/config"
	"shortener/internal/middleware"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler POST /api/auth/login
func LoginHandler(c config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 校验凭证
		if req.Username != c.Auth.Admin.Username || req.Password != c.Auth.Admin.Password {
			w.WriteHeader(http.StatusUnauthorized)
			httpx.OkJsonCtx(r.Context(), w, map[string]string{"msg": "用户名或密码错误"})
			return
		}

		// 生成 JWT
		token, err := middleware.GenerateToken(req.Username, c.Auth.JWTSecret, c.Auth.Expiry)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 设置 HTTP-only Cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			MaxAge:   c.Auth.Expiry,
			HttpOnly: true,
			Secure:   true, // 生产环境 HTTPS 必须为 true
			SameSite: http.SameSiteLaxMode,
		})

		httpx.OkJsonCtx(r.Context(), w, map[string]string{"message": "登录成功"})
	}
}

// LogoutHandler POST /api/auth/logout
func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		httpx.OkJsonCtx(r.Context(), w, map[string]string{"message": "已登出"})
	}
}

// MeHandler GET /api/auth/me
func MeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(middleware.UserKey).(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]string{"username": username})
	}
}
