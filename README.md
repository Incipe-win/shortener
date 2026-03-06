# Smart-Shortener Gateway

基于微服务与大模型的智能短链安全网关，采用 [go-zero](https://github.com/zeromicro/go-zero) 微服务框架构建。

## ✨ 核心亮点

### 🚀 高性能架构
- 基于 go-zero 微服务框架构建，采用 **Base62 算法 + 发号器模式** 实现高效长短链映射
- 运用 **Bloom Filter**（20MB，容纳千万级数据）前置拦截非法访问，彻底解决恶意伪造短链导致的缓存穿透问题
- **Redis 缓存** 加速短链查询，命中率显著提升

### 🤖 AI 内容理解与赋能
- 集成大语言模型（LLM），在长链录入时 **异步提取页面特征**
- 支持生成基于语义的 **个性化可读短链**（如 `go-tutorial`、`python-data`）
- 提供 **AI 目标页面摘要预览**，提升终端用户点击转化率

### 🔒 零信任安全机制
- 转链前引入 **AI 风险评估 + 自动化安全巡检**（域名黑名单 → URL 特征分析 → LLM 深度评估）
- 跳转时根据风险等级 **分级处理**：`danger` 拒绝跳转 / `warning` 附加警告 / `safe` 正常 302
- **黑名单双重拦截**：保留词拦截（`api`、`admin` 等）+ 恶意域名拦截

### 📊 工程化与可观测性
- 集成 **OpenTelemetry** 全链路追踪，关键路径自定义 Span
- 集成 **Prometheus** 自定义 Metrics（`/metrics` 端点），涵盖转链/跳转/BloomFilter/LLM 延迟/安全拦截等指标
- **防循环嵌套检测**：域名层面 (`IsCircularURL`) + URL BasePath 双重拦截

## 📁 项目结构

```
shortener/
├── shortener.go              # 程序入口
├── shortener.api             # API 定义文件
├── etc/
│   └── shortener-api.yaml    # 配置文件（gitignore, 用 .example 模板）
├── internal/
│   ├── config/config.go      # 配置结构体
│   ├── handler/              # HTTP 处理器
│   │   ├── routes.go         # 路由注册
│   │   ├── convertHandler.go # 转链处理
│   │   ├── showHandler.go    # 跳转处理
│   │   └── previewHandler.go # AI 预览处理
│   ├── logic/                # 业务逻辑
│   │   ├── convertLogic.go   # 转链逻辑（含异步 AI 分析）
│   │   ├── showLogic.go      # 跳转逻辑（含安全过滤）
│   │   └── previewLogic.go   # AI 预览逻辑
│   ├── svc/serviceContext.go # 服务上下文
│   └── types/types.go        # 请求/响应类型
├── model/                    # 数据库 Model 层
├── pkg/                      # 工具包
│   ├── base62/               # Base62 编解码
│   ├── connect/              # URL 连通性检查
│   ├── llm/                  # LLM 客户端封装
│   ├── md5/                  # MD5 工具
│   ├── metrics/              # Prometheus 指标定义
│   ├── otel/                 # OpenTelemetry 初始化
│   ├── safety/               # URL 安全巡检
│   ├── scraper/              # 页面内容抓取
│   └── urltool/              # URL 工具（防循环检测）
├── short_url_map.sql         # 建表 DDL
├── sequence.sql              # 发号器表 DDL
└── alter_short_url_map.sql   # AI 字段升级 DDL
```

## 🛠 快速开始

### 1. 环境依赖

- Go 1.24+
- MySQL 5.7+
- Redis 6.0+

### 2. 建库建表

```sql
CREATE DATABASE IF NOT EXISTS sql_test DEFAULT CHARSET utf8mb4;
USE sql_test;
```

执行以下 SQL 文件：
```bash
mysql -u root -p sql_test < sequence.sql
mysql -u root -p sql_test < short_url_map.sql
```

> 如果是从旧版本升级，只需执行 `alter_short_url_map.sql` 添加 AI 字段。

### 3. 修改配置

复制配置模板：
```bash
cp etc/shortener-api.yaml.example etc/shortener-api.yaml
```

根据实际环境修改 `etc/shortener-api.yaml` 中的数据库和 Redis 连接信息。

### 4. 安装依赖 & 启动

```bash
go mod tidy
go run .
```

看到输出即表示启动成功：
```
Starting server at 0.0.0.0:8888...
```

### 5. 启用 AI 功能（可选）

在 `shortener-api.yaml` 中设置：
```yaml
LLM:
  Enabled: true
```

启动时通过环境变量注入 API Key（**不要将密钥写入配置文件**）：
```bash
LLM_API_KEY="sk-your-key" go run .
```

## 📡 API 接口

### 转链 — `POST /convert`

将长链接转换为短链接，同时异步触发 AI 分析。

```bash
curl -X POST http://localhost:8888/convert \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://www.example.com"}'
```

响应：
```json
{ "short_url": "shortener.com/1En" }
```

### 跳转 — `GET /:short_url`

302 跳转到原始长链接（danger 级别会被拦截）。

```bash
curl -L http://localhost:8888/1En
```

### AI 预览 — `GET /preview/:short_url`

获取短链对应页面的 AI 摘要、关键词和安全等级。

```bash
curl http://localhost:8888/preview/1En
```

响应：
```json
{
  "short_url": "1En",
  "long_url": "https://www.example.com",
  "summary": "这是一个示例网站...",
  "keywords": ["示例", "网站"],
  "risk_level": "safe"
}
```

### Prometheus 指标 — `GET /metrics`

暴露 Prometheus 格式的监控指标。

## ⚙️ 配置说明

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `ShortUrlDB.DSN` | MySQL 连接串 | — |
| `Sequence.DSN` | 发号器表连接串 | — |
| `CacheRedis` | Redis 缓存配置 | — |
| `BaseString` | Base62 编码字符表（打乱顺序提高安全性） | — |
| `ShortUrlBlackList` | 短链保留词黑名单 | `["api","convert"...]` |
| `ShortDoamin` | 短域名 | `shortener.com` |
| `LLM.Enabled` | 是否启用 AI 功能 | `false` |
| `LLM.BaseURL` | LLM API 地址（OpenAI 兼容） | — |
| `LLM.APIKey` | API 密钥（建议用 `${ENV_VAR}` 注入） | — |
| `LLM.Model` | 模型名称 | `gpt-4o-mini` |
| `Safety.Enabled` | 是否启用安全巡检 | `true` |
| `Safety.BlackListDomains` | 恶意域名黑名单 | `[]` |
| `Otel.Name` | 服务名称 | `shortener-api` |
| `Otel.Endpoint` | OTLP Collector 地址（留空则不启用） | `""` |
| `Otel.Sampler` | 采样率 (0.0~1.0) | `1.0` |

## 🧪 运行测试

```bash
go test ./pkg/... -v
```

## 📖 goctl 生成命令参考

```bash
# 根据 api 文件生成代码
goctl api go -api shortener.api -dir . -style=goZero

# 根据数据表生成 model 层代码
goctl model mysql datasource -url="root:root@tcp(127.0.0.1:3306)/sql_test" -table="short_url_map" -dir="./model" -c
goctl model mysql datasource -url="root:root@tcp(127.0.0.1:3306)/sql_test" -table="sequence" -dir="./model" -c
```

> ⚠️ 重新生成 model 后需手动补回 `FindAll()`、`UpdateAIFields()` 等自定义方法，以及 `ShortUrlMap` 结构体中的 AI 字段。

## 参考

- [Go 语言微服务和云原生课程](https://github.com/Q1mi/go-micro-service-and-cloud-native-course)
- [go-zero 官方文档](https://go-zero.dev/)
