package handler

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	"shortener/internal/config"
	"shortener/internal/middleware"
	"shortener/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler POST /api/auth/login
func LoginHandler(c config.Config, userModel model.UserModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 查找用户
		user, err := userModel.FindOneByUsername(r.Context(), req.Username)
		if err == sqlx.ErrNotFound {
			w.WriteHeader(http.StatusUnauthorized)
			httpx.OkJsonCtx(r.Context(), w, map[string]string{"msg": "用户名或密码错误"})
			return
		}
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 校验密码
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			httpx.OkJsonCtx(r.Context(), w, map[string]string{"msg": "用户名或密码错误"})
			return
		}

		// 生成 JWT
		token, err := middleware.GenerateToken(user.Username, user.Id, c.Auth.JWTSecret, c.Auth.Expiry)
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
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})

		httpx.OkJsonCtx(r.Context(), w, map[string]any{"message": "登录成功", "username": user.Username, "user_id": user.Id})
	}
}

// RegisterHandler POST /api/auth/register
func RegisterHandler(c config.Config, userModel model.UserModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req registerRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 校验用户名：3-32 字符，仅字母数字
		if matched, _ := regexp.MatchString(`^[a-zA-Z0-9]{3,32}$`, req.Username); !matched {
			httpx.Error(w, errors.New("用户名需 3-32 位字母或数字"))
			return
		}

		// 校验密码：至少 6 字符
		if len(req.Password) < 6 {
			httpx.Error(w, errors.New("密码至少 6 位"))
			return
		}

		// 检查用户名是否已存在
		_, err := userModel.FindOneByUsername(r.Context(), req.Username)
		if err != sqlx.ErrNotFound {
			if err == nil {
				httpx.Error(w, errors.New("用户名已存在"))
				return
			}
			logx.Errorw("failed to check username", logx.LogField{Key: "err", Value: err.Error()})
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 哈希密码
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 插入用户
		res, err := userModel.Insert(r.Context(), &model.User{
			Username:     req.Username,
			PasswordHash: string(hash),
		})
		if err != nil {
			logx.Errorw("failed to insert user", logx.LogField{Key: "err", Value: err.Error()})
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		userID, err := res.LastInsertId()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 注册成功，自动登录
		token, err := middleware.GenerateToken(req.Username, uint64(userID), c.Auth.JWTSecret, c.Auth.Expiry)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			MaxAge:   c.Auth.Expiry,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})

		httpx.OkJsonCtx(r.Context(), w, map[string]any{"message": "注册成功", "username": req.Username, "user_id": userID})
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
		username, userID, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]any{"username": username, "user_id": userID})
	}
}
