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
		Offset: int32(offset),
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
		utils.NewError(w, http.StatusBadRequest, "could not find Content-Type")
		return
	}

	if withImage {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			utils.NewError(w, http.StatusNotFound, "could not parse form")
			return
		}
	}
	title := strings.TrimSpace(r.FormValue("title"))
	text := strings.TrimSpace(r.FormValue("text"))
	userId, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		utils.NewError(w, http.StatusInternalServerError, err.Error())
		return
	}
	suboliveId, err := strconv.Atoi(r.FormValue("subolive_id"))
	if err != nil {
		utils.NewError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if withImage {
		image, header, err := r.FormFile("image")
		if err != nil {
			utils.NewError(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer image.Close()

		errs, valid, imgPath:= utils.ValidateNewPostWithImage(title, text, image, header)
		if !valid {
			utils.NewErrorBody(w, http.StatusInternalServerError, consts.ResCreatedPost{
				Post: sqlc.Post{
					PostID: 0,
					Title: title,
					Text: text,
					CreatedAt: time.Now(),
					UserID: 0,
					SuboliveID: 0,
					ImageID: sql.NullInt32{
						Int32: 0,
						Valid: false,
					},
				},
				Errors: errs,
			})
			return
		}
		errors = errs

		_, err = h.q.CreateImage(context.Background(), imgPath)
		if err != nil {
			utils.NewError(w, http.StatusInternalServerError, err.Error())
			return
		}

		img, err := h.q.GetNewestImage(context.Background())
		if err != nil {
			utils.NewError(w, http.StatusInternalServerError, err.Error())
			return
		}

		_, err = h.q.CreatePost(context.Background(), sqlc.CreatePostParams{
			Title: title,
			Text: text,
			UserID: int32(userId),
			SuboliveID: int32(suboliveId),
			ImageID: sql.NullInt32{
				Int32: img.ImageID,
				Valid: true,
			},
		})

	} else {
		errs, valid := utils.ValidateNewPost(title, text)
		if !valid {
			utils.NewErrorBody(w, http.StatusInternalServerError, consts.ResCreatedPost{
				Post: sqlc.Post{
					PostID: 0,
					Title: title,
					Text: text,
					CreatedAt: time.Now(),
					UserID: 0,
					SuboliveID: 0,
					ImageID: sql.NullInt32{
						Int32: 0,
						Valid: false,
					},
				},
				Errors: errs,
			})
			return
		}
		errors = errs

		_, err = h.q.CreatePost(context.Background(), sqlc.CreatePostParams{
			Title: title,
			Text: text,
			UserID: int32(userId),
			SuboliveID: int32(suboliveId),
			ImageID: sql.NullInt32{
				Int32: 0,
				Valid: false,
			},
		})
	}

	post, err := h.q.GetNewestPost(context.Background())
	if err != nil {
		utils.NewError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resPost := sqlc.Post{
		PostID: post.PostID,
		Title: post.Title,
		Text: post.Text,
		CreatedAt: post.CreatedAt,
		UserID: post.UserID,
		SuboliveID: post.SuboliveID,
		ImageID: post.ImageID,
	}

	res := consts.ResCreatedPost{
		Post: resPost,
		Errors: errors,
	}

	utils.NewResponse(w, http.StatusCreated, res)
	return
}















