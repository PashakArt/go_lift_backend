SELECT user_id, tenant_id, telegram_id, phone, created_at
FROM users
WHERE tenant_id = $1 AND telegram_id = $2;