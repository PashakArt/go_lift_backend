SELECT 
    e.exercise_id, 
    e.name, 
    e.type, 
    e.is_global
FROM exercises e
JOIN exercise_muscle_groups emg ON e.exercise_id = emg.exercise_id
LEFT JOIN user_favorite_exercises f 
    ON e.exercise_id = f.exercise_id AND f.user_id = $2
WHERE 
    emg.muscle_group_id = $1
    AND e.is_global = true
ORDER BY 
    f.user_id DESC NULLS LAST,
    e.name ASC;