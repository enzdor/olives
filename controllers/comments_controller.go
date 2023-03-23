package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jobutterfly/olives/consts"
	"github.com/jobutterfly/olives/sqlc"
	"github.com/jobutterfly/olives/utils"
)

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.NewResponse(w, http.StatusMethodNotAllowed, consts.ResCreatedComment{
			Comment:       consts.EmptyComment,
			FormErrors: consts.EmptyCreateCommentErrors,
			Error:      consts.UnsupportedMethod.Error(),
		})
		return
	}

	contentType := r.Header.Get("Content-Type")
	content := strings.Split(contentType, ";")
	var withImage bool = false
	var errors [2]consts.FormInputError
	switch content[0] {
	case "application/x-www-form-urlencoded":
		withImage = false
	case "multipart/form-data":
		withImage = true
	default:
		utils.NewResponse(w, http.StatusBadRequest, consts.ResCreatedComment{
			Comment:       consts.EmptyComment,
			FormErrors: consts.EmptyCreateCommentErrors,
			Error:      "Counld not find Content-Type",
		})
		return
	}

	if withImage {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedComment{
				Comment:       consts.EmptyComment,
				FormErrors: consts.EmptyCreateCommentErrors,
				Error:      "Could not parse form",
			})
			return
		}
	}
	text := strings.TrimSpace(r.FormValue("text"))
	userId, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedComment{
			Comment:       consts.EmptyComment,
			FormErrors: consts.EmptyCreateCommentErrors,
			Error:      err.Error(),
		})
		return
	}
	postId, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedComment{
			Comment:       consts.EmptyComment,
			FormErrors: consts.EmptyCreateCommentErrors,
			Error:      err.Error(),
		})
		return
	}

	if withImage {
		image, header, err := r.FormFile("image")
		if err != nil {
			utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedComment{
				Comment:       consts.EmptyComment,
				FormErrors: consts.EmptyCreateCommentErrors,
				Error:      err.Error(),
			})
			return
		}
		defer image.Close()

		errs, valid, imgPath := utils.ValidateNewCommentWithImage(text, image, header)
		if !valid {
			utils.NewResponse(w, http.StatusUnprocessableEntity, consts.ResCreatedComment{
				Comment: sqlc.Comment{
					CommentID: 0,
					Text:      text,
					CreatedAt: time.Now(),
					UserID:    int32(userId),
					ImageID: sql.NullInt32{
						Int32: 0,
						Valid: false,
					},
					PostID: int32(postId),
				},
				FormErrors: errs,
				Error:      "",
			})
			return
		}
		errors = errs

		_, err = h.q.CreateImage(context.Background(), imgPath)
		if err != nil {
			utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedComment{
				Comment:       consts.EmptyComment,
				FormErrors: consts.EmptyCreateCommentErrors,
				Error:      err.Error(),
			})
			return
		}

		img, err := h.q.GetNewestImage(context.Background())
		if err != nil {
			utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedComment{
				Comment:       consts.EmptyComment,
				FormErrors: consts.EmptyCreateCommentErrors,
				Error:      err.Error(),
			})
			return
		}

		_, err = h.q.CreateComment(context.Background(), sqlc.CreateCommentParams{
			Text:   text,
			UserID: int32(userId),
			ImageID: sql.NullInt32{
				Int32: img.ImageID,
				Valid: true,
			},
			PostID: int32(postId),
		})

	} else {
		errs, valid := utils.ValidateNewComment(text)
		if !valid {
			utils.NewResponse(w, http.StatusUnprocessableEntity, consts.ResCreatedComment{
				Comment: sqlc.Comment{
					CommentID: 0,
					Text:      text,
					CreatedAt: time.Now(),
					UserID:    int32(userId),
					ImageID: sql.NullInt32{
						Int32: 0,
						Valid: false,
					},
					PostID: int32(postId),
				},
				FormErrors: errs,
				Error:      "",
			})
			return
		}
		errors = errs

		_, err = h.q.CreateComment(context.Background(), sqlc.CreateCommentParams{
			Text:   text,
			UserID: int32(userId),
			ImageID: sql.NullInt32{
				Int32: 0,
				Valid: false,
			},
			PostID: int32(postId),
		})
	}

	comment, err := h.q.GetNewestComment(context.Background())
	if err != nil {
		utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedComment{
			Comment:       consts.EmptyComment,
			FormErrors: consts.EmptyCreateCommentErrors,
			Error:      err.Error(),
		})
		return
	}

	resComment := sqlc.Comment{
		CommentID: comment.CommentID,
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt,
		UserID:    comment.UserID,
		ImageID:   comment.ImageID,
		PostID:    comment.PostID,
	}

	res := consts.ResCreatedComment{
		Comment:       resComment,
		FormErrors: errors,
		Error:      "",
	}

	utils.NewResponse(w, http.StatusCreated, res)
	return
}
