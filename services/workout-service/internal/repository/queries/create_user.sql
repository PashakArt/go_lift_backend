INSERT INTO users (user_id, tenant_id, telegram_id, phone, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING user_id, created_at;