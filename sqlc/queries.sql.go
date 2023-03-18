// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: queries.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"
)

const countComments = `-- name: CountComments :one
SELECT COUNT(*) FROM comments
WHERE post_id = ?
`

func (q *Queries) CountComments(ctx context.Context, postID int32) (int64, error) {
	row := q.db.QueryRowContext(ctx, countComments, postID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countPosts = `-- name: CountPosts :one
SELECT COUNT(*) FROM posts
`

func (q *Queries) CountPosts(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countPosts)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createComment = `-- name: CreateComment :execresult
INSERT INTO posts(text, created_at, user_id, image_id, post_id)
VALUES (?, ?, ?, ?, ?)
`

type CreateCommentParams struct {
	Text      string        `json:"text"`
	CreatedAt time.Time     `json:"created_at"`
	UserID    int32         `json:"user_id"`
	ImageID   sql.NullInt32 `json:"image_id"`
	PostID    int32         `json:"post_id"`
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createComment,
		arg.Text,
		arg.CreatedAt,
		arg.UserID,
		arg.ImageID,
		arg.PostID,
	)
}

const createImage = `-- name: CreateImage :execresult
INSERT INTO images(file_path)
VALUES(?)
`

func (q *Queries) CreateImage(ctx context.Context, filePath string) (sql.Result, error) {
	return q.db.ExecContext(ctx, createImage, filePath)
}

const createPost = `-- name: CreatePost :execresult
INSERT INTO posts(title, text, user_id, image_id, subolive_id)
VALUES (?, ?, ?, ?, ?)
`

type CreatePostParams struct {
	Title      string        `json:"title"`
	Text       string        `json:"text"`
	UserID     int32         `json:"user_id"`
	ImageID    sql.NullInt32 `json:"image_id"`
	SuboliveID int32         `json:"subolive_id"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createPost,
		arg.Title,
		arg.Text,
		arg.UserID,
		arg.ImageID,
		arg.SuboliveID,
	)
}

const createSubolive = `-- name: CreateSubolive :execresult
INSERT INTO subolives(name)
VALUES(?)
`

func (q *Queries) CreateSubolive(ctx context.Context, name string) (sql.Result, error) {
	return q.db.ExecContext(ctx, createSubolive, name)
}

const createUser = `-- name: CreateUser :execresult
INSERT INTO users (email, username, password)
VALUES (?, ?, ?)
`

type CreateUserParams struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createUser, arg.Email, arg.Username, arg.Password)
}

const deletePost = `-- name: DeletePost :execresult
DELETE FROM posts
WHERE post_id = ?
`

func (q *Queries) DeletePost(ctx context.Context, postID int32) (sql.Result, error) {
	return q.db.ExecContext(ctx, deletePost, postID)
}

const deleteUser = `-- name: DeleteUser :execresult
DELETE FROM users
WHERE user_id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, userID int32) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteUser, userID)
}

const getNewestImage = `-- name: GetNewestImage :one
SELECT image_id, file_path FROM images
WHERE image_id = (
	SELECT MAX(image_id) FROM images
)
`

func (q *Queries) GetNewestImage(ctx context.Context) (Image, error) {
	row := q.db.QueryRowContext(ctx, getNewestImage)
	var i Image
	err := row.Scan(&i.ImageID, &i.FilePath)
	return i, err
}

const getNewestPost = `-- name: GetNewestPost :one
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
	users.email, 
	users.password 
FROM posts
LEFT JOIN subolives ON posts.subolive_id = subolives.subolive_id
LEFT JOIN users ON posts.user_id = users.user_id
LEFT JOIN images ON posts.image_id = images.image_id
ORDER BY created_at DESC, post_id DESC
LIMIT 1
`

type GetNewestPostRow struct {
	PostID     int32          `json:"post_id"`
	Title      string         `json:"title"`
	Text       string         `json:"text"`
	CreatedAt  time.Time      `json:"created_at"`
	SuboliveID int32          `json:"subolive_id"`
	Name       sql.NullString `json:"name"`
	ImageID    sql.NullInt32  `json:"image_id"`
	FilePath   sql.NullString `json:"file_path"`
	UserID     int32          `json:"user_id"`
	Username   sql.NullString `json:"username"`
	Email      sql.NullString `json:"email"`
	Password   sql.NullString `json:"password"`
}

func (q *Queries) GetNewestPost(ctx context.Context) (GetNewestPostRow, error) {
	row := q.db.QueryRowContext(ctx, getNewestPost)
	var i GetNewestPostRow
	err := row.Scan(
		&i.PostID,
		&i.Title,
		&i.Text,
		&i.CreatedAt,
		&i.SuboliveID,
		&i.Name,
		&i.ImageID,
		&i.FilePath,
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.Password,
	)
	return i, err
}

const getNewestUser = `-- name: GetNewestUser :one
SELECT user_id, email, username, password FROM users
WHERE user_id = (
	SELECT MAX(user_id) FROM users
)
`

func (q *Queries) GetNewestUser(ctx context.Context) (User, error) {
	row := q.db.QueryRowContext(ctx, getNewestUser)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Username,
		&i.Password,
	)
	return i, err
}

const getPost = `-- name: GetPost :one
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
	users.email, 
	users.password 
FROM posts
LEFT JOIN subolives ON posts.subolive_id = subolives.subolive_id
LEFT JOIN users ON posts.user_id = users.user_id
LEFT JOIN images ON posts.image_id = images.image_id
WHERE post_id = ?
LIMIT 1
`

type GetPostRow struct {
	PostID     int32          `json:"post_id"`
	Title      string         `json:"title"`
	Text       string         `json:"text"`
	CreatedAt  time.Time      `json:"created_at"`
	SuboliveID int32          `json:"subolive_id"`
	Name       sql.NullString `json:"name"`
	ImageID    sql.NullInt32  `json:"image_id"`
	FilePath   sql.NullString `json:"file_path"`
	UserID     int32          `json:"user_id"`
	Username   sql.NullString `json:"username"`
	Email      sql.NullString `json:"email"`
	Password   sql.NullString `json:"password"`
}

func (q *Queries) GetPost(ctx context.Context, postID int32) (GetPostRow, error) {
	row := q.db.QueryRowContext(ctx, getPost, postID)
	var i GetPostRow
	err := row.Scan(
		&i.PostID,
		&i.Title,
		&i.Text,
		&i.CreatedAt,
		&i.SuboliveID,
		&i.Name,
		&i.ImageID,
		&i.FilePath,
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.Password,
	)
	return i, err
}

const getPosts = `-- name: GetPosts :many
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
	users.email, 
	users.password 
FROM posts
LEFT JOIN subolives ON posts.subolive_id = subolives.subolive_id
LEFT JOIN users ON posts.user_id = users.user_id
LEFT JOIN images ON posts.image_id = images.image_id
ORDER BY created_at ASC
LIMIT ?
`

type GetPostsRow struct {
	PostID     int32          `json:"post_id"`
	Title      string         `json:"title"`
	Text       string         `json:"text"`
	CreatedAt  time.Time      `json:"created_at"`
	SuboliveID int32          `json:"subolive_id"`
	Name       sql.NullString `json:"name"`
	ImageID    sql.NullInt32  `json:"image_id"`
	FilePath   sql.NullString `json:"file_path"`
	UserID     int32          `json:"user_id"`
	Username   sql.NullString `json:"username"`
	Email      sql.NullString `json:"email"`
	Password   sql.NullString `json:"password"`
}

func (q *Queries) GetPosts(ctx context.Context, limit int32) ([]GetPostsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPosts, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsRow
	for rows.Next() {
		var i GetPostsRow
		if err := rows.Scan(
			&i.PostID,
			&i.Title,
			&i.Text,
			&i.CreatedAt,
			&i.SuboliveID,
			&i.Name,
			&i.ImageID,
			&i.FilePath,
			&i.UserID,
			&i.Username,
			&i.Email,
			&i.Password,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSubolivePosts = `-- name: GetSubolivePosts :many
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
	users.email, 
	users.password 
FROM posts
LEFT JOIN subolives ON posts.subolive_id = subolives.subolive_id
LEFT JOIN users ON posts.user_id = users.user_id
LEFT JOIN images ON posts.image_id = images.image_id
WHERE posts.subolive_id = ?
ORDER BY created_at ASC
LIMIT 10
OFFSET ?
`

type GetSubolivePostsParams struct {
	SuboliveID int32 `json:"subolive_id"`
	Offset     int32 `json:"offset"`
}

type GetSubolivePostsRow struct {
	PostID     int32          `json:"post_id"`
	Title      string         `json:"title"`
	Text       string         `json:"text"`
	CreatedAt  time.Time      `json:"created_at"`
	SuboliveID int32          `json:"subolive_id"`
	Name       sql.NullString `json:"name"`
	ImageID    sql.NullInt32  `json:"image_id"`
	FilePath   sql.NullString `json:"file_path"`
	UserID     int32          `json:"user_id"`
	Username   sql.NullString `json:"username"`
	Email      sql.NullString `json:"email"`
	Password   sql.NullString `json:"password"`
}

func (q *Queries) GetSubolivePosts(ctx context.Context, arg GetSubolivePostsParams) ([]GetSubolivePostsRow, error) {
	rows, err := q.db.QueryContext(ctx, getSubolivePosts, arg.SuboliveID, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSubolivePostsRow
	for rows.Next() {
		var i GetSubolivePostsRow
		if err := rows.Scan(
			&i.PostID,
			&i.Title,
			&i.Text,
			&i.CreatedAt,
			&i.SuboliveID,
			&i.Name,
			&i.ImageID,
			&i.FilePath,
			&i.UserID,
			&i.Username,
			&i.Email,
			&i.Password,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSubolives = `-- name: GetSubolives :many
SELECT subolive_id, name FROM subolives
`

func (q *Queries) GetSubolives(ctx context.Context) ([]Subolife, error) {
	rows, err := q.db.QueryContext(ctx, getSubolives)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Subolife
	for rows.Next() {
		var i Subolife
		if err := rows.Scan(&i.SuboliveID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
SELECT user_id, email, username, password FROM users
WHERE user_id = ?
LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, userID int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Username,
		&i.Password,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT user_id, email, username, password FROM users
WHERE email = ?
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Username,
		&i.Password,
	)
	return i, err
}
