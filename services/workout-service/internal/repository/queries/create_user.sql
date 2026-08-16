INSERT INTO users (user_id, tenant_id, telegram_id, phone, created_at, tg_username, tg_first_name, tg_last_name)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING user_id, created_at;