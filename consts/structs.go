package consts

import (
	"database/sql"
	"time"

	"github.com/jobutterfly/olives/sqlc"
)

var EmptyPost = sqlc.Post{
	PostID:     0,
	Title:      "",
	Text:       "",
	CreatedAt:  time.Now(),
	UserID:     0,
	SuboliveID: 0,
	ImageID: sql.NullInt32{
		Int32: 0,
		Valid: false,
	},
}

var EmptyUser = sqlc.User{
	UserID:   0,
	Email:    "",
	Username: "",
	Password: "",
	Admin:    false,
}

var EmptyComment = sqlc.Comment{
	CommentID: 0,
	Text:      "",
	CreatedAt: time.Now(),
	UserID:    0,
	ImageID: sql.NullInt32{
		Int32: 0,
		Valid: false,
	},
	PostID: 0,
}
