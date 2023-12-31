-- name: GetPosts :many
SELECT posts.post_id, 
	posts.title, 
	posts.text, 
	posts.created_at, 
	posts.subolive_id, 
	subolives.name, 
	posts.image_id, 
	images.file_path, 
	posts.user_id, 
	users.username, 
	users.email
FROM posts
LEFT JOIN subolives ON posts.subolive_id = subolives.subolive_id
LEFT JOIN users ON posts.user_id = users.user_id
LEFT JOIN images ON posts.image_id = images.image_id
ORDER BY created_at ASC
LIMIT ?;

-- name: GetPost :one
SELECT posts.post_id, 
	posts.title, 
	posts.text, 
	posts.created_at, 
	posts.subolive_id, 
	subolives.name, 
	posts.image_id, 
	images.file_path, 
	posts.user_id, 
	users.username, 
	users.email
FROM posts
LEFT JOIN subolives ON posts.subolive_id = subolives.subolive_id
LEFT JOIN users ON posts.user_id = users.user_id
LEFT JOIN images ON posts.image_id = images.image_id
WHERE post_id = ?
LIMIT 1;

-- name: GetNewestPost :one
SELECT posts.post_id, 
	posts.title, 
	posts.text, 
	posts.created_at, 
	posts.subolive_id, 
	subolives.name, 
	posts.image_id, 
	images.file_path, 
	posts.user_id, 
	users.username, 
	users.email
FROM posts
LEFT JOIN subolives ON posts.subolive_id = subolives.subolive_id
LEFT JOIN users ON posts.user_id = users.user_id
LEFT JOIN images ON posts.image_id = images.image_id
ORDER BY created_at DESC, post_id DESC
LIMIT 1;

-- name: GetSubolivePosts :many
SELECT posts.post_id, 
	posts.title, 
	posts.text, 
	posts.created_at, 
	posts.subolive_id, 
	subolives.name, 
	posts.image_id, 
	images.file_path, 
	posts.user_id, 
	users.username, 
	users.email
FROM posts
LEFT JOIN subolives ON posts.subolive_id = subolives.subolive_id
LEFT JOIN users ON posts.user_id = users.user_id
LEFT JOIN images ON posts.image_id = images.image_id
WHERE posts.subolive_id = ?
ORDER BY created_at ASC
LIMIT 10
OFFSET ?;

-- name: CreatePost :execresult
INSERT INTO posts(title, text, user_id, image_id, subolive_id)
VALUES (?, ?, ?, ?, ?);

-- name: DeletePost :execresult
DELETE FROM posts
WHERE post_id = ?;

-- name: CountPosts :one
SELECT COUNT(*) FROM posts;

-- name: GetComment :one
SELECT comments.comment_id, 
	comments.text, 
	comments.created_at, 
	comments.post_id, 
	comments.image_id, 
	images.file_path, 
	comments.user_id, 
	users.username, 
	users.email
FROM comments
LEFT JOIN posts ON comments.post_id = posts.post_id
LEFT JOIN users ON comments.user_id = users.user_id
LEFT JOIN images ON comments.image_id = images.image_id
WHERE comment_id = ?
LIMIT 1;

-- name: CreateComment :execresult
INSERT INTO comments(text, created_at, user_id, image_id, post_id)
VALUES (?, ?, ?, ?, ?);

-- name: DeleteComment :execresult
DELETE FROM comments
WHERE comment_id = ?;

-- name: CountComments :one
SELECT COUNT(*) FROM comments;

-- name: GetNewestComment :one
SELECT * FROM comments
WHERE comment_id = (
	SELECT MAX(comment_id) FROM comments
);

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = ?
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?
LIMIT 1;

-- name: GetNewestUser :one
SELECT * FROM users
WHERE user_id = (
	SELECT MAX(user_id) FROM users
);

-- name: CreateUser :execresult
INSERT INTO users (email, username, password, admin)
VALUES (?, ?, ?, ?);

-- name: DeleteUser :execresult
DELETE FROM users
WHERE user_id = ?;

-- name: CreateImage :execresult
INSERT INTO images(file_path)
VALUES(?);

-- name: GetNewestImage :one
SELECT * FROM images
WHERE image_id = (
	SELECT MAX(image_id) FROM images
);

-- name: DeleteImage :execresult
DELETE FROM images
WHERE image_id = ?;

-- name: GetSubolives :many
SELECT * FROM subolives;

-- name: CreateSubolive :execresult
INSERT INTO subolives(name)
VALUES(?);

-- name: GetSession :one
SELECT * FROM sessions
LEFT JOIN users ON sessions.user_id = users.user_id
WHERE sessions.session_id = ?;

-- name: GetNewestSession :one
SELECT * FROM sessions
WHERE session_id = (
	SELECT MAX(session_id) FROM sessions
);

-- name: CreateSession :execresult
INSERT INTO sessions(session_id, user_id)
VALUES(?, ?);

-- name: UpdateSession :execresult
UPDATE sessions
SET last_access = CURRENT_TIMESTAMP()
WHERE session_id = ?;

-- name: DeleteSession :execresult
DELETE FROM sessions
WHERE session_id = ?;






