SELECT tenant_id, branding_json
FROM tenants
WHERE tenant_id = $1