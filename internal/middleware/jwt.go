package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type contextKey string

const UserKey contextKey = "auth_user"

// Claims JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
func GenerateToken(username, secret string, expiry int) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiry) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "shortener-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// JWTCookieMiddleware 从 HTTP-only Cookie 中读取 JWT 并验证
func JWTCookieMiddleware(secret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				httpx.Error(w, &httpError{Code: http.StatusUnauthorized, Msg: "未登录"})
				return
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				httpx.Error(w, &httpError{Code: http.StatusUnauthorized, Msg: "登录已过期"})
				return
			}

			// 将用户信息注入 context
			ctx := context.WithValue(r.Context(), UserKey, claims.Username)
			next(w, r.WithContext(ctx))
		}
	}
}

// httpError 用于返回标准 HTTP 错误
type httpError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *httpError) Error() string {
	return e.Msg
}
