BACKUP_DIR="backups"
mkdir -p "$BACKUP_DIR"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
docker exec dotapro-postgres pg_dump \
  --username=postgres \
  --dbname=dotapro \
  --no-owner \
  --no-acl \
  --format=plain \
  > "$BACKUP_DIR/dotapro_${TIMESTAMP}.sql"

echo "✅ Backup created: $BACKUP_DIR/dotapro_${TIMESTAMP}.sql"