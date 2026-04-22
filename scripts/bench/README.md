# Smart-Shortener 压力测试

## 快速开始

### 1. 安装 wrk

```bash
# Ubuntu/Debian
sudo apt install -y wrk

# macOS
brew install wrk
```

### 2. 启动监控栈（Prometheus + Grafana）

```bash
# 先启动基础服务
docker compose up -d

# 再启动监控组件
docker compose -f docker-compose.yaml -f docker-compose.bench.yaml up -d
```

启动后访问：
- **Grafana**: http://localhost:3001 (admin/admin) → 自动加载 "Smart-Shortener 压测看板"
- **Prometheus**: http://localhost:9090 → 查询指标
- **pprof**: http://localhost:8888/debug/pprof/

### 3. 运行压测

```bash
# 默认压测（4 线程，100 并发，30 秒）
./scripts/bench/run.sh http://localhost:8888

# 自定义压测参数
BENCH_DURATION=60s BENCH_CONNECTIONS=500 BENCH_THREADS=8 \
  ./scripts/bench/run.sh http://localhost:8888
```

### 4. pprof 性能分析

```bash
# CPU 分析（30 秒采样）
curl -o cpu.prof http://localhost:8888/debug/pprof/profile?seconds=30
go tool pprof -http=:6060 cpu.prof

# 内存分析
curl -o mem.prof http://localhost:8888/debug/pprof/heap
go tool pprof -http=:6061 mem.prof

# 阻塞分析
curl -o block.prof http://localhost:8888/debug/pprof/block?seconds=30
go tool pprof -http=:6062 block.prof

# 运行 go tool trace（需要先生成 trace）
curl -o trace.out http://localhost:8888/debug/pprof/trace?seconds=10
go tool trace trace.out
```

### 5. 单独 wrk 压测

```bash
# 仅转链接口
wrk -t4 -c100 -d30s -s scripts/bench/convert.lua http://localhost:8888

# 仅跳转接口（需先创建短 URL）
SHORT_URLS="abc123,def456" wrk -t4 -c100 -d30s -s scripts/bench/show.lua http://localhost:8888

# 仅预览接口
SHORT_URLS="abc123,def456" wrk -t4 -c100 -d30s -s scripts/bench/preview.lua http://localhost:8888
```

## 压测策略

### 阶梯式 QPS 测试

逐步增加并发连接数，找到系统瓶颈：

```bash
for conn in 10 50 100 200 500 1000; do
  echo "=== Connections: $conn ==="
  wrk -t4 -c$conn -d10s -s scripts/bench/convert.lua http://localhost:8888
done
```

### 稳定性压测

长时间运行验证内存泄漏：

```bash
wrk -t4 -c100 -d10m -s scripts/bench/convert.lua http://localhost:8888
```

同时观察 `go tool pprof heap` 的内存增长趋势。

## 关键指标解读

| 指标 | 说明 |
|------|------|
| QPS | 每秒请求数，越高越好 |
| P50 | 中位数延迟，反映正常用户体验 |
| P99 | 99% 请求的延迟上限，反映长尾问题 |
| Bloom Filter hit rate | 应接近 100%（已存在的 URL） |
| Kafka 生产/消费速率 | 应基本平衡，积压说明消费跟不上 |
