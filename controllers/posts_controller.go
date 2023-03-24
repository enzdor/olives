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

func (h *Handler) GetOrDeletePost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		h.DeletePost(w, r)
		return
	case "GET":
		h.GetPost(w, r)
		return
	default:
		utils.NewError(w, http.StatusMethodNotAllowed, consts.UnsupportedMethod.Error())
	}

}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	v, err := utils.GetPathValues(strings.Split(r.URL.Path, "/"), 0)
	if err != nil {
		utils.NewError(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.q.GetPost(context.Background(), int32(v.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			utils.NewError(w, http.StatusNotFound, err.Error())
			return
		}

		utils.NewError(w, http.StatusInternalServerError, "error when getting post")
		return
	}

	utils.NewResponse(w, http.StatusOK, post)
	return

}

func (h *Handler) GetSubolivePosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.NewError(w, http.StatusMethodNotAllowed, consts.UnsupportedMethod.Error())
		return
	}

	v, err := utils.GetPathValues(strings.Split(r.URL.Path, "/"), 1)
	if err != nil {
		utils.NewError(w, http.StatusBadRequest, err.Error())
		return
	}

	if v.Id > 5 {
		utils.NewError(w, http.StatusNotFound, consts.SuboliveNonExistant.Error())
		return
	}

	queryPage := r.URL.Query().Get("page")
	var page int
	if queryPage == "" {
		page = 0
	} else {
		intPage, err := strconv.Atoi(queryPage)
		if err != nil {
			utils.NewError(w, http.StatusInternalServerError, consts.PageNotAnInteger.Error())
			return
		}
		page = intPage
	}
	offset := page * consts.ITEMS_PER_PAGE

	posts, err := h.q.GetSubolivePosts(context.Background(), sqlc.GetSubolivePostsParams{
		Offset:     int32(offset),
		SuboliveID: int32(v.Id),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			utils.NewError(w, http.StatusNotFound, err.Error())
			return
		}

		utils.NewError(w, http.StatusInternalServerError, "error when getting posts")
		return
	}

	if len(posts) < 1 {
		utils.NewError(w, http.StatusNotFound, sql.ErrNoRows.Error())
		return
	}

	utils.NewResponse(w, http.StatusOK, posts)
	return
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.NewResponse(w, http.StatusMethodNotAllowed, consts.ResCreatedPost{
			Post:       consts.EmptyPost,
			FormErrors: consts.EmptyCreatePostErrors,
			Error:      consts.UnsupportedMethod.Error(),
		})
		return
	}

	contentType := r.Header.Get("Content-Type")
	content := strings.Split(contentType, ";")
	var withImage bool = false
	var errors [3]consts.FormInputError
	switch content[0] {
	case "application/x-www-form-urlencoded":
		withImage = false
	case "multipart/form-data":
		withImage = true
	default:
		utils.NewResponse(w, http.StatusBadRequest, consts.ResCreatedPost{
			Post:       consts.EmptyPost,
			FormErrors: consts.EmptyCreatePostErrors,
			Error:      "Counld not find Content-Type",
		})
		return
	}

	if withImage {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedPost{
				Post:       consts.EmptyPost,
				FormErrors: consts.EmptyCreatePostErrors,
				Error:      "Could not parse form",
			})
			return
		}
	}
	title := strings.TrimSpace(r.FormValue("title"))
	text := strings.TrimSpace(r.FormValue("text"))
	userId, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedPost{
			Post:       consts.EmptyPost,
			FormErrors: consts.EmptyCreatePostErrors,
			Error:      err.Error(),
		})
		return
	}
	suboliveId, err := strconv.Atoi(r.FormValue("subolive_id"))
	if err != nil {
		utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedPost{
			Post:       consts.EmptyPost,
			FormErrors: consts.EmptyCreatePostErrors,
			Error:      err.Error(),
		})
		return
	}

	if withImage {
		image, header, err := r.FormFile("image")
		if err != nil {
			utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedPost{
				Post:       consts.EmptyPost,
				FormErrors: consts.EmptyCreatePostErrors,
				Error:      err.Error(),
			})
			return
		}
		defer image.Close()

		errs, valid, imgPath := utils.ValidateNewPostWithImage(title, text, image, header)
		if !valid {
			utils.NewResponse(w, http.StatusUnprocessableEntity, consts.ResCreatedPost{
				Post: sqlc.Post{
					PostID:     0,
					Title:      title,
					Text:       text,
					CreatedAt:  time.Now(),
					UserID:     int32(userId),
					SuboliveID: int32(suboliveId),
					ImageID: sql.NullInt32{
						Int32: 0,
						Valid: false,
					},
				},
				FormErrors: errs,
				Error:      "",
			})
			return
		}
		errors = errs

		_, err = h.q.CreateImage(context.Background(), imgPath)
		if err != nil {
			utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedPost{
				Post:       consts.EmptyPost,
				FormErrors: consts.EmptyCreatePostErrors,
				Error:      err.Error(),
			})
			return
		}

		img, err := h.q.GetNewestImage(context.Background())
		if err != nil {
			utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedPost{
				Post:       consts.EmptyPost,
				FormErrors: consts.EmptyCreatePostErrors,
				Error:      err.Error(),
			})
			return
		}

		_, err = h.q.CreatePost(context.Background(), sqlc.CreatePostParams{
			Title:      title,
			Text:       text,
			UserID:     int32(userId),
			SuboliveID: int32(suboliveId),
			ImageID: sql.NullInt32{
				Int32: img.ImageID,
				Valid: true,
			},
		})

	} else {
		errs, valid := utils.ValidateNewPost(title, text)
		if !valid {
			utils.NewResponse(w, http.StatusUnprocessableEntity, consts.ResCreatedPost{
				Post: sqlc.Post{
					PostID:     0,
					Title:      title,
					Text:       text,
					CreatedAt:  time.Now(),
					UserID:     int32(userId),
					SuboliveID: int32(suboliveId),
					ImageID: sql.NullInt32{
						Int32: 0,
						Valid: false,
					},
				},
				FormErrors: errs,
				Error:      "",
			})
			return
		}
		errors = errs

		_, err = h.q.CreatePost(context.Background(), sqlc.CreatePostParams{
			Title:      title,
			Text:       text,
			UserID:     int32(userId),
			SuboliveID: int32(suboliveId),
			ImageID: sql.NullInt32{
				Int32: 0,
				Valid: false,
			},
		})
	}

	post, err := h.q.GetNewestPost(context.Background())
	if err != nil {
		utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreatedPost{
			Post:       consts.EmptyPost,
			FormErrors: consts.EmptyCreatePostErrors,
			Error:      err.Error(),
		})
		return
	}

	resPost := sqlc.Post{
		PostID:     post.PostID,
		Title:      post.Title,
		Text:       post.Text,
		CreatedAt:  post.CreatedAt,
		UserID:     post.UserID,
		SuboliveID: post.SuboliveID,
		ImageID:    post.ImageID,
	}

	res := consts.ResCreatedPost{
		Post:       resPost,
		FormErrors: errors,
		Error:      "",
	}

	utils.NewResponse(w, http.StatusCreated, res)
	return
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	v, err := utils.GetPathValues(strings.Split(r.URL.Path, "/"), 0)
	if err != nil {
		utils.NewError(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.q.GetPost(context.Background(), int32(v.Id))
	if post.FilePath.Valid {
		if err := os.Remove(post.FilePath.String); err != nil {
			utils.NewError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	exc, err := h.q.DeletePost(context.Background(), int32(v.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			utils.NewError(w, http.StatusNotFound, err.Error())
			return
		}

		utils.NewError(w, http.StatusInternalServerError, "error when deleting post")
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

	utils.NewResponse(w, http.StatusOK, "")
}
