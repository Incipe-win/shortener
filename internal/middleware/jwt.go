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
const UserIDKey contextKey = "auth_user_id"

// Claims JWT claims
type Claims struct {
	Username string `json:"username"`
	UserID   uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
func GenerateToken(username string, userID uint64, secret string, expiry int) (string, error) {
	claims := Claims{
		Username: username,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiry) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "shortener-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// JWTCookieMiddleware 从 HTTP-only Cookie 中读取 JWT 并验证（必须认证）
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
			ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
			next(w, r.WithContext(ctx))
		}
	}
}

// OptionalJWTCookieMiddleware 从 HTTP-only Cookie 中读取 JWT（可选认证）
// 如果没有 token 或 token 无效，继续执行但注入空用户信息
func OptionalJWTCookieMiddleware(secret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				// 没有 token，注入空用户信息继续
				ctx := context.WithValue(r.Context(), UserKey, "")
				ctx = context.WithValue(ctx, UserIDKey, uint64(0))
				next(w, r.WithContext(ctx))
				return
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				// token 无效，注入空用户信息继续
				ctx := context.WithValue(r.Context(), UserKey, "")
				ctx = context.WithValue(ctx, UserIDKey, uint64(0))
				next(w, r.WithContext(ctx))
				return
			}

			// token 有效，注入用户信息
			ctx := context.WithValue(r.Context(), UserKey, claims.Username)
			ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
			next(w, r.WithContext(ctx))
		}
	}
}

// GetUserFromContext 从 context 中获取用户信息
func GetUserFromContext(ctx context.Context) (username string, userID uint64, ok bool) {
	u, ok1 := ctx.Value(UserKey).(string)
	id, ok2 := ctx.Value(UserIDKey).(uint64)
	return u, id, ok1 && ok2
}

// httpError 用于返回标准 HTTP 错误
type httpError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *httpError) Error() string {
	return e.Msg
}
