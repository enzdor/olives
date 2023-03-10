-- name: GetPosts :many
SELECT * FROM posts
LEFT JOIN subolives ON posts.subolive_id = subolives.subolive_id
LEFT JOIN users ON posts.user_id = users.user_id
LEFT JOIN images ON posts.image_id = images.image_id
ORDER BY created_at ASC
LIMIT ?;

-- name: GetPost :one
SELECT * FROM posts
LEFT JOIN subolives ON posts.subolive_id = subolives.subolive_id
LEFT JOIN users ON posts.user_id = users.user_id
LEFT JOIN images ON posts.image_id = images.image_id
WHERE post_id = ?
LIMIT 1;

-- name: GetSubolivePosts :many
SELECT * FROM posts
LEFT JOIN subolives ON posts.subolive_id = subolives.subolive_id
LEFT JOIN users ON posts.user_id = users.user_id
LEFT JOIN images ON posts.image_id = images.image_id
WHERE posts.subolive_id = ?
ORDER BY created_at ASC
LIMIT 10
OFFSET ?;

-- name: CreatePost :execresult
INSERT INTO posts(title, text, created_at, user_id, image_id, subolive_id)
VALUES (?, ?, ?, ?, ?, ?);

-- name: DeletePost :execresult
DELETE FROM posts
WHERE post_id = ?;

-- name: CountPosts :one
SELECT COUNT(*) FROM posts;

-- name: CreateComment :execresult
INSERT INTO posts(text, created_at, user_id, image_id, post_id)
VALUES (?, ?, ?, ?, ?);

-- name: CountComments :one
SELECT COUNT(*) FROM comments
WHERE post_id = ?;

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = ?
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?
LIMIT 1;

-- name: CreateUser :execresult
INSERT INTO users (email, username, password)
VALUES (?, ?, ?);

-- name: DeleteUser :execresult
DELETE FROM users
WHERE user_id = ?;

-- name: CreateImage :execresult
INSERT INTO images(file_path)
VALUES(?);

-- name: GetSubolives :many
SELECT * FROM subolives;

-- name: CreateSubolive :execresult
INSERT INTO subolives(name)
VALUES(?);









