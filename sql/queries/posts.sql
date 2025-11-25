-- name: CreatePost :one

INSERT INTO posts (id, body, user_id, created_at, updated_at) 
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    NOW()
)
RETURNING *;


-- name: GetPost :one

SELECT * FROM posts
WHERE id = $1;

-- name: GetPosts :many

SELECT * FROM posts 
ORDER BY created_at;

-- name: DeletePost :one

DELETE  FROM posts 
WHERE id = $1
RETURNING *;

-- name: GetPostsByUserId :many
SELECT * FROM posts 
WHERE user_id = $1
ORDER BY created_at;