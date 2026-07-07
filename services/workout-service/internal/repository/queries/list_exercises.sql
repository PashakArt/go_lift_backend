SELECT 
    e.exercise_id, 
    e.tenant_id, 
    e.name, 
    e.type, 
    e.is_global, 
    e.created_at,
    COALESCE(array_agg(mg.code) FILTER (WHERE mg.code IS NOT NULL), '{}') as muscle_group_codes
FROM exercises e
LEFT JOIN exercise_muscle_groups emg ON e.exercise_id = emg.exercise_id
LEFT JOIN muscle_groups mg ON emg.muscle_group_id = mg.muscle_group_id
WHERE (e.is_global = true OR e.tenant_id = $1)
GROUP BY e.exercise_id, e.tenant_id, e.name, e.type, e.is_global, e.created_at
-- Фильтрация по конкретной мышце (если передан пустой $2, условие игнорируется)
HAVING $2 = '' OR $2 = ANY(COALESCE(array_agg(mg.code) FILTER (WHERE mg.code IS NOT NULL), '{}'))
ORDER BY e.name ASC;