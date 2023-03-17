package controllers

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		utils.NewError(w, http.StatusInternalServerError, "could not parse file")
		return
	}

	img, _, err := r.FormFile("image")
	if err != nil {
		utils.NewError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer img.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, img); err != nil {
		utils.NewError(w, http.StatusInternalServerError, "could not read image")
		return
	}
	
	if err := os.WriteFile("created.png", buf.Bytes(), 0666); err != nil {
		utils.NewError(w, http.StatusInternalServerError, "could not write image")
		return
	}
}















