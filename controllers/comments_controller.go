package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jobutterfly/olives/consts"
	"github.com/jobutterfly/olives/sqlc"
	"github.com/jobutterfly/olives/utils"
)

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if err := h.Authorizer(r, false); err != nil {
		utils.NewError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if r.Method != http.MethodPost {
		utils.NewResponse(w, http.StatusMethodNotAllowed, consts.ResCreatedComment{
			Comment:    consts.EmptyComment,
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
			Comment:    consts.EmptyComment,
			FormErrors: consts.EmptyCreateCommentErrors,
			Error:      "Counld not find Content-Type",
		})
		return
	}

	if withImage {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedComment{
				Comment:    consts.EmptyComment,
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
			Comment:    consts.EmptyComment,
			FormErrors: consts.EmptyCreateCommentErrors,
			Error:      err.Error(),
		})
		return
	}
	postId, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedComment{
			Comment:    consts.EmptyComment,
			FormErrors: consts.EmptyCreateCommentErrors,
			Error:      err.Error(),
		})
		return
	}

	if withImage {
		image, header, err := r.FormFile("image")
		if err != nil {
			utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedComment{
				Comment:    consts.EmptyComment,
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
				Comment:    consts.EmptyComment,
				FormErrors: consts.EmptyCreateCommentErrors,
				Error:      err.Error(),
			})
			return
		}

		img, err := h.q.GetNewestImage(context.Background())
		if err != nil {
			utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedComment{
				Comment:    consts.EmptyComment,
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
			Comment:    consts.EmptyComment,
			FormErrors: consts.EmptyCreateCommentErrors,
			Error:      err.Error(),
		})
		return
	}

	res := consts.ResCreatedComment{
		Comment:    comment,
		FormErrors: errors,
		Error:      "",
	}

	utils.NewResponse(w, http.StatusCreated, res)
	return
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if err := h.Authorizer(r, true); err != nil {
		utils.NewError(w, http.StatusUnauthorized, err.Error())
		return
	}

	v, err := utils.GetPathValues(strings.Split(r.URL.Path, "/"), 0)
	if err != nil {
		utils.NewError(w, http.StatusBadRequest, err.Error())
		return
	}

	comment, err := h.q.GetComment(context.Background(), int32(v.Id))

	exc, err := h.q.DeleteComment(context.Background(), int32(v.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			utils.NewError(w, http.StatusNotFound, err.Error())
			return
		}

		utils.NewError(w, http.StatusInternalServerError, "error when deleting comment")
		return
	}

	rows, err := exc.RowsAffected()
	if err != nil {
		utils.NewError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if rows < 1 {
		utils.NewError(w, http.StatusNotFound, sql.ErrNoRows.Error())
		return
	}

	if comment.FilePath.Valid {
		if err := os.Remove(comment.FilePath.String); err != nil {
			utils.NewError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if comment.ImageID.Valid {
			if _, err := h.q.DeleteImage(context.Background(), comment.ImageID.Int32); err != nil {
				utils.NewError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	utils.NewResponse(w, http.StatusOK, "")
}
