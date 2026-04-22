#!/bin/bash
set -euo pipefail

# Smart-Shortener 压测脚本
# 用法: ./scripts/bench/run.sh [TARGET]
#   TARGET: 后端地址，默认 http://localhost:8888

TARGET="${1:-http://localhost:8888}"
DURATION="${BENCH_DURATION:-30s}"
CONNECTIONS="${BENCH_CONNECTIONS:-100}"
THREADS="${BENCH_THREADS:-4}"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}=== Smart-Shortener 压力测试 ===${NC}"
echo "Target: $TARGET"
echo "Duration: $DURATION | Connections: $CONNECTIONS | Threads: $THREADS"
echo ""

# 检查 wrk 是否安装
if ! command -v wrk &>/dev/null; then
    echo -e "${RED}wrk 未安装，请先安装:${NC}"
    echo "  Ubuntu/Debian: sudo apt install -y wrk"
    echo "  macOS:         brew install wrk"
    exit 1
fi

# ── Step 1: 检查后端健康 ──
echo -e "${YELLOW}[1/4] 健康检查...${NC}"
if curl -sf "$TARGET/health" > /dev/null 2>&1; then
    echo "  后端运行中"
else
    echo -e "${RED}  后端不可达: $TARGET${NC}"
    exit 1
fi

# ── Step 2: 预热 — 创建一批短 URL 供 show/preview 测试使用 ──
echo -e "${YELLOW}[2/4] 预热：创建短 URL...${NC}"
SHORT_URL_LIST=""
for i in $(seq 1 20); do
    resp=$(curl -sf -X POST "$TARGET/api/convert" \
        -H 'Content-Type: application/json' \
        -d '{"long_url":"https://example.com/page'"$i"'"}' 2>/dev/null || true)
    if [ -n "$resp" ]; then
        surl=$(echo "$resp" | grep -oP '(?<=/)[A-Za-z0-9]+$' || true)
        if [ -n "$surl" ]; then
            SHORT_URL_LIST="${SHORT_URL_LIST:+$SHORT_URL_LIST,}$surl"
        fi
    fi
done

if [ -z "$SHORT_URL_LIST" ]; then
    echo -e "${RED}  预热失败，无法创建短 URL，可能已达到未注册用户限制${NC}"
    echo "  跳过 show/preview 压测，仅执行 convert 压测"
fi
echo "  已创建 $(echo "$SHORT_URL_LIST" | tr ',' '\n' | grep -c . || echo 0) 个短 URL 用于压测"

# 创建输出目录
BENCH_DIR="$(dirname "$0")/results"
mkdir -p "$BENCH_DIR"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# ── Step 3: 压测 /api/convert ──
echo -e "${YELLOW}[3/4] 压测 POST /api/convert...${NC}"
wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" \
    -s "$(dirname "$0")/convert.lua" "$TARGET/api/convert" \
    -o "$BENCH_DIR/convert_${TIMESTAMP}.txt" 2>&1 | tee "$BENCH_DIR/convert_${TIMESTAMP}.log"
echo ""

# ── Step 4: 压测 /:short_url (如果有预热数据) ──
if [ -n "$SHORT_URL_LIST" ]; then
    echo -e "${YELLOW}[4/4] 压测 GET /:short_url (redirect)...${NC}"
    SHORT_URLS="$SHORT_URL_LIST" \
    wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" \
        -s "$(dirname "$0")/show.lua" "$TARGET" \
        -o "$BENCH_DIR/show_${TIMESTAMP}.txt" 2>&1 | tee "$BENCH_DIR/show_${TIMESTAMP}.log"
    echo ""

    echo -e "${YELLOW}  压测 GET /api/preview/:short_url...${NC}"
    SHORT_URLS="$SHORT_URL_LIST" \
    wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" \
        -s "$(dirname "$0")/preview.lua" "$TARGET/api/preview" \
        -o "$BENCH_DIR/preview_${TIMESTAMP}.txt" 2>&1 | tee "$BENCH_DIR/preview_${TIMESTAMP}.log"
else
    echo -e "${YELLOW}[4/4] 跳过 show/preview 压测（无预热数据）${NC}"
fi

echo ""
echo -e "${GREEN}=== 压测完成 ===${NC}"
echo "结果保存在: $BENCH_DIR/"

# ── 快速汇总 ──
echo ""
echo -e "${YELLOW}--- QPS 快速汇总 ---${NC}"
for log in "$BENCH_DIR"/*_${TIMESTAMP}.log; do
    [ -f "$log" ] || continue
    name=$(basename "$log" .log)
    qps=$(grep -oP 'Requests/sec:\s+\K[0-9.]+' "$log" || echo "N/A")
    p99=$(grep -oP '99th percentile\s+\K[0-9.]+[a-z]+' "$log" || echo "N/A")
    echo "  $name: QPS=$qps | P99=$p99"
done
