package ctxdata

// context key 类型，避免冲突
type CtxKey string

const (
	// KeyClientIP 客户端 IP
	KeyClientIP CtxKey = "client_ip"
	// KeyUserAgent User-Agent
	KeyUserAgent CtxKey = "user_agent"
	// KeyReferer Referer
	KeyReferer CtxKey = "referer"
)
