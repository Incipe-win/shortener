# ── 阶段1: 编译 ──────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /build

# 安装 ca-certificates（生产环境 HTTPS 需要）
RUN apk add --no-cache git ca-certificates

# 先复制 go mod 文件，利用 Docker 层缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o /build/server .

# ── 阶段2: 运行 ──────────────────────────
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -u 1000 appuser

WORKDIR /app

# 从编译阶段复制二进制和配置文件
COPY --from=builder /build/server .
COPY --from=builder /build/etc ./etc

USER appuser

EXPOSE 8888

HEALTHCHECK --interval=10s --timeout=3s --start-period=30s --retries=3 \
  CMD wget -qO- http://localhost:8888/health || exit 1

ENTRYPOINT ["./server"]
CMD ["-f", "etc/shortener-api.yaml"]
