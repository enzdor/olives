package consts

type CreatePostRequest struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	SuboliveID int32  `json:"subolive_id"`
	UserID     int32  `json:"user_id"`
	ImageID    string `json:"image_id"`
}
