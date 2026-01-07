
SELECT
task_tags.id as tag_id,
task_tags.name as tag_name,
task_templates.id,
task_templates.name,
task_templates.description,
task_templates.created_at,
task_templates.updated_at,
task_templates.deleted_at
FROM task_templates
LEFT JOIN task_template_task_tags as pivot on task_templates.id = pivot.task_template_id
LEFT JOIN task_tags on pivot.task_tag_id = task_tags.id
WHERE task_templates.id = ?

SELECT
task_tags.id as tag_id,
task_tags.name as tag_name,
task_templates.id,
task_templates.name,
task_templates.description,
task_templates.created_at,
task_templates.updated_at,
task_templates.deleted_at
FROM task_templates
LEFT JOIN task_template_task_tags as pivot on task_templates.id = pivot.task_template_id
LEFT JOIN task_tags on pivot.task_tag_id = task_tags.id
WHERE task_templates.id = 1;
