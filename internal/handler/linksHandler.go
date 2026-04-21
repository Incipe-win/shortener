package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"shortener/internal/middleware"
	"shortener/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type linkItem struct {
	ID         uint64   `json:"id"`
	Surl       string   `json:"surl"`
	Lurl       string   `json:"lurl"`
	AISummary  string   `json:"ai_summary"`
	AIKeywords []string `json:"ai_keywords"`
	RiskLevel  string   `json:"risk_level"`
	ClickCount int64    `json:"click_count"`
	CreateAt   string   `json:"create_at"`
}

type linksResponse struct {
	List     []linkItem `json:"list"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// LinksHandler GET /api/links?page=1&page_size=10&search=xxx
func LinksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, userID, _ := middleware.GetUserFromContext(r.Context())

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
		if pageSize < 1 || pageSize > 100 {
			pageSize = 10
		}
		search := r.URL.Query().Get("search")

		// 查询总数（按用户过滤）
		total, err := svcCtx.ShortUrlModel.Count(r.Context(), userID, search)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 查询列表（按用户过滤）
		list, err := svcCtx.ShortUrlModel.FindList(r.Context(), userID, page, pageSize, search)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 转换为响应结构
		items := make([]linkItem, 0, len(list))
		for _, m := range list {
			item := linkItem{
				ID:       m.Id,
				CreateAt: m.CreateAt.Format("2006-01-02T15:04:05Z07:00"),
				ClickCount: int64(m.ClickCount),
			}
			if m.Surl.Valid {
				item.Surl = m.Surl.String
			}
			if m.Lurl.Valid {
				item.Lurl = m.Lurl.String
			}
			if m.AiSummary.Valid {
				item.AISummary = m.AiSummary.String
			}
			if m.AiKeywords.Valid {
				_ = json.Unmarshal([]byte(m.AiKeywords.String), &item.AIKeywords)
			}
			if m.RiskLevel.Valid {
				item.RiskLevel = m.RiskLevel.String
			}
			items = append(items, item)
		}

		httpx.OkJsonCtx(r.Context(), w, linksResponse{
			List:     items,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		})
	}
}
