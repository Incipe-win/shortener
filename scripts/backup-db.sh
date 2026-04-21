#!/usr/bin/env bash
# ── 数据库备份脚本 ──────────────────────────
# 用法: ./scripts/backup-db.sh
# 保留最近 7 天的备份
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKUP_DIR="${SCRIPT_DIR}/../backups"
DATE="$(date +%Y%m%d_%H%M%S)"
RETENTION_DAYS=7

# 从 .env 读取配置
ENV_FILE="${SCRIPT_DIR}/../.env"
if [ -f "$ENV_FILE" ]; then
  set -a
  source "$ENV_FILE"
  set +a
fi

DB_USER="${MYSQL_USER:-shortener}"
DB_PASS="${MYSQL_PASSWORD:-}"
DB_NAME="${MYSQL_DATABASE:-shortener}"

mkdir -p "$BACKUP_DIR"

echo "[$(date)] 开始备份数据库 ${DB_NAME}..."
docker exec shortener-mysql mysqldump \
  -u"$DB_USER" \
  -p"$DB_PASS" \
  --single-transaction \
  --routines \
  --triggers \
  --databases "$DB_NAME" \
  | gzip > "${BACKUP_DIR}/${DB_NAME}_${DATE}.sql.gz"

SIZE=$(du -h "${BACKUP_DIR}/${DB_NAME}_${DATE}.sql.gz" | cut -f1)
echo "[$(date)] 备份完成: ${BACKUP_DIR}/${DB_NAME}_${DATE}.sql.gz (${SIZE})"

# 清理过期备份
echo "[$(date)] 清理 ${RETENTION_DAYS} 天前的备份..."
find "$BACKUP_DIR" -name "${DB_NAME}_*.sql.gz" -mtime +${RETENTION_DAYS} -delete
echo "[$(date)] 完成"
