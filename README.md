# Smart-Shortener Gateway

基于微服务与大模型的智能短链安全网关。

## 架构

```
用户 → incipe.top:2053 (HTTPS)
  → Nginx (TLS 终止 + 限流)
    → /api/* → Go 后端 (:8888)
    → /:short_url → Go 后端 (302 重定向)
    → /* → React 前端 SPA
```

## 技术栈

- 后端: Go 1.24 + go-zero 微服务框架
- 前端: React 19 + TypeScript + TailwindCSS V4
- 数据库: MySQL 8.4
- 缓存: Redis 7 + Bloom Filter
- 消息队列: Apache Kafka
- 链路追踪: OpenTelemetry + Jaeger
- 部署: Docker Compose + Nginx 反向代理

## 功能

- 短链接生成与跳转 (base62 编码)
- 用户注册与 JWT Cookie 认证
- 多用户数据隔离 — 每个用户仅可见自己的链接与统计数据
- 未注册用户限制创建 3 个短链接，注册后无限制
- 安全巡检与风险评级 (LLM AI 分析)
- 链接预览 (摘要 + 关键词提取)
- 基于 Redis 的 IP 限流
- 点击事件异步统计
- Prometheus 指标采集

## 快速开始

### 前置条件

- Docker + Docker Compose v2
- 域名 TLS 证书 (cert.pem + key.pem)

### 1. 配置环境变量

```bash
cp .env.example .env
```

编辑 `.env`，填写所有密码和密钥：

```
MYSQL_ROOT_PASSWORD=your-secure-root-password
MYSQL_PASSWORD=your-secure-db-password
REDIS_PASSWORD=your-secure-redis-password
JWT_SECRET=your-super-secret-jwt-key
LLM_API_KEY=your-deepseek-api-key
DOMAIN=incipe.top
HTTPS_PORT=2053
```

### 2. 放置 TLS 证书

```bash
mkdir -p nginx/ssl
cp your-cert.pem nginx/ssl/cert.pem
cp your-key.pem nginx/ssl/key.pem
```

### 3. 构建并启动

```bash
make docker-build
make up
```

服务启动后访问: `https://incipe.top:2053`

### 4. 验证

```bash
# 健康检查
curl -k https://localhost:2053/health

# 创建短链接（未注册用户，最多 3 次）
curl -k https://localhost:2053/api/convert \
  -H 'Content-Type: application/json' \
  -d '{"long_url":"https://example.com"}'

# 注册账号
curl -k https://localhost:2053/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"myuser","password":"mypass123"}'

# 短链接跳转
curl -kI https://localhost:2053/abc123
```

### 5. 运维命令

```bash
make logs       # 查看所有服务日志
make down       # 停止所有服务
make restart    # 重启所有服务
make backup     # 数据库备份
make clean      # 清理 Docker 资源
```

## 压力测试与性能分析

### 安装压测工具

```bash
# Ubuntu/Debian
sudo apt install -y wrk

# macOS
brew install wrk
```

### 启动监控栈（Prometheus + Grafana）

```bash
docker compose -f docker-compose.yaml -f docker-compose.bench.yaml up -d
```

启动后访问：
- **Grafana**: http://localhost:3001 (admin/admin) — 自动加载压测看板
- **Prometheus**: http://localhost:9090 — 查询指标
- **pprof**: http://localhost:8888/debug/pprof/ — 性能分析

### 一键压测

```bash
# 默认压测（4 线程，100 并发，30 秒）
./scripts/bench/run.sh http://localhost:8888

# 自定义参数
BENCH_DURATION=60s BENCH_CONNECTIONS=500 BENCH_THREADS=8 \
  ./scripts/bench/run.sh http://localhost:8888
```

### pprof 性能分析

```bash
# CPU 热点分析（30 秒采样）
curl -o cpu.prof http://localhost:8888/debug/pprof/profile?seconds=30
go tool pprof -http=:6060 cpu.prof

# 内存分配分析
curl -o mem.prof http://localhost:8888/debug/pprof/heap
go tool pprof -http=:6060 mem.prof

# 完整调用链追踪
curl -o trace.out http://localhost:8888/debug/pprof/trace?seconds=10
go tool trace trace.out
```

### 阶梯式 QPS 测试

逐步增加并发连接数，找到系统瓶颈：

```bash
for conn in 10 50 100 200 500 1000; do
  echo "=== Connections: $conn ==="
  wrk -t4 -c$conn -d10s -s scripts/bench/convert.lua http://localhost:8888
done
```

> 详见 `scripts/bench/README.md`

## 本地开发

### 后端

```bash
# 1. 启动基础设施 (仅 MySQL + Redis + Kafka + Jaeger)
docker compose up -d mysql redis kafka jaeger

# 2. 运行后端
go run shortener.go

# 3. 运行前端
cd web && pnpm install && pnpm dev
```

后端默认监听 `0.0.0.0:8888`，前端开发服务器监听 `0.0.0.0:3000` (已配置 API 代理)。

### 前端

```bash
cd web
pnpm install
pnpm dev        # 开发模式 (热更新 + API 代理)
pnpm build      # 生产构建
```

## 项目结构

```
├── Dockerfile                # Go 后端多阶段构建
├── Makefile                  # 统一构建/部署命令
├── docker-compose.yaml       # 全栈服务编排
├── .env.example              # 环境变量模板
├── etc/
│   └── shortener-api.yaml    # go-zero 配置文件
├── internal/
│   ├── config/               # 配置结构体
│   ├── consumer/             # Kafka 消费者 (AI/点击/安全)
│   ├── handler/              # HTTP 处理器
│   ├── logic/                # 业务逻辑
│   ├── middleware/           # 中间件 (JWT/CORS/限流)
│   ├── svc/                  # 服务上下文
│   └── types/                # 类型定义
├── model/                    # MySQL Model 层 (goctl 生成)
├── nginx/
│   ├── nginx.conf            # Nginx 主配置 (公网入口)
│   └── shortener.conf        # TLS + 反向代理配置
│   └── ssl/                  # TLS 证书目录 (git ignored)
├── pkg/                      # 公共包 (LLM/Kafka/OTel)
├── scripts/
│   └── backup-db.sh          # 数据库备份脚本
├── shortener.api             # API 定义 (goctl)
├── shortener.go              # 主入口
└── web/
    ├── Dockerfile            # 前端多阶段构建
    ├── nginx.conf            # 前端容器 Nginx 配置
    ├── src/                  # React 前端源码
    └── package.json
```

## API 接口

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/convert` | 创建短链接（未注册用户限 3 次） |
| GET | `/api/convert/remaining` | 查询未注册用户剩余可创建次数 |
| GET | `/api/preview/:short_url` | 链接预览 (AI 摘要 + 风险评级) |
| GET | `/:short_url` | 短链接跳转 (302) |

### 认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/auth/register` | 用户注册 (注册后自动登录) |
| POST | `/api/auth/login` | 用户登录 |
| POST | `/api/auth/logout` | 登出 |
| GET | `/api/auth/me` | 获取当前用户信息 |
| GET | `/api/links?page=1&page_size=10` | 链接列表 (分页 + 搜索，仅当前用户) |
| GET | `/api/stats` | 仪表盘统计 (仅当前用户) |
| GET | `/api/metrics` | Prometheus 指标 |

## 数据库

初始化 SQL 文件会在首次启动时自动执行：

- `short_url_map.sql` -- 短链接映射表
- `sequence.sql` -- 发号器表
- `alter_short_url_map.sql` -- 表结构扩展
- `alter_click_count.sql` -- 点击计数字段
- `alter_user.sql` -- 用户表 + 用户 ID 字段

数据持久化到 Docker volume `mysql_data`。

## 用户与权限

| 用户类型 | 短链接创建 | 数据可见性 |
|----------|-----------|-----------|
| 未注册 | 最多 3 个 | 仅创建时返回结果，无后台 |
| 已注册 | 无限制 | 仅自己的链接和统计数据 |

- 注册时自动登录，JWT Cookie 有效期默认 7 天
- 用户名 3-32 位字母数字，密码至少 6 位
- 密码使用 bcrypt 哈希存储

## 备份

```bash
make backup
```

备份文件存储在 `backups/` 目录，保留最近 7 天。可通过 crontab 定时执行：

```
0 2 * * * cd /path/to/shortener && make backup
```
