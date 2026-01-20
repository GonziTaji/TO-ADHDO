WITH input_ids(tag_id) AS (
    SELECT 'ropa'
    UNION ALL SELECT 'asdf'
)
SELECT
at.tag_id
FROM articles_tags at
WHERE at.article_id = ?
AND at.deleted_at IS NULL
AND at.tag_id NOT IN (
    SELECT tag_id
    FROM input_ids
);





WITH input_ids(tag_id) AS (
    SELECT 1
    UNION ALL SELECT 2
    UNION ALL SELECT 4
)
SELECT
at.tag_id,
CASE WHEN ii.tag_id IS NULL THEN 0 ELSE 1 END as in_input,
    CASE WHEN at.tag_id IS NULL THEN 0 ELSE 1 END as in_relation
        FROM articles_tags at
        LEFT JOIN input_ids as ii
        ON at.article_id = ii
        WHERE at.article_id = 11
        AND at.deleted_at IS NULL;


