package consts

import "time"

type Post struct {
	PostID     int32     `json:"post_id"`
	Title      string    `json:"title"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
	SuboliveID int32     `json:"subolive_id"`
	Name       string    `json:"name"`
	ImageID    int32     `json:"image_id"`
	FilePath   string    `json:"file_path"`
	UserID     int32     `json:"user_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
}
