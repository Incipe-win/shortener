.PHONY: help build docker-build up down logs restart clean backup

# 默认目标
help:
	@echo "Smart-Shortener 部署命令"
	@echo ""
	@echo "  make build         编译 Go 后端 + 构建前端"
	@echo "  make docker-build  构建所有 Docker 镜像"
	@echo "  make up            启动所有服务"
	@echo "  make down          停止所有服务"
	@echo "  make restart       重启所有服务"
	@echo "  make logs          查看服务日志"
	@echo "  make backup        数据库备份"
	@echo "  make clean         清理 Docker 资源"

# ── 本地构建 ──────────────────────────────────
build:
	@echo "编译 Go 后端..."
	CGO_ENABLED=0 go build -ldflags='-s -w' -o bin/server .
	@echo "构建前端..."
	cd web && pnpm build
	@echo "构建完成"

# ── Docker 构建 ───────────────────────────────
docker-build:
	@echo "构建所有 Docker 镜像..."
	docker compose build --pull
	@echo "构建完成"

# ── 启动服务 ──────────────────────────────────
up:
	@echo "启动服务..."
	docker compose up -d
	@echo "等待服务就绪..."
	@sleep 5
	@docker compose ps
	@echo ""
	@echo "访问地址: https://$$(grep '^DOMAIN=' .env 2>/dev/null | cut -d= -f2 || echo 'incipe.top'):$$(grep '^HTTPS_PORT=' .env 2>/dev/null | cut -d= -f2 || echo '2053')"

# ── 停止服务 ──────────────────────────────────
down:
	docker compose down

# ── 重启服务 ──────────────────────────────────
restart: down up

# ── 查看日志 ──────────────────────────────────
logs:
	docker compose logs -f --tail=100

# ── 数据库备份 ────────────────────────────────
backup:
	./scripts/backup-db.sh

# ── 清理 Docker 资源 ──────────────────────────
clean:
	docker compose down -v --remove-orphans
	docker system prune -f
