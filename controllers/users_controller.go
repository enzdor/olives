package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/jobutterfly/olives/consts"
	"github.com/jobutterfly/olives/sqlc"
	"github.com/jobutterfly/olives/utils"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) GetOrDeleteUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		h.DeleteUser(w, r)
		return
	case "GET":
		h.GetUser(w, r)
		return
	default:
		utils.NewError(w, http.StatusMethodNotAllowed, consts.UnsupportedMethod.Error())
	}

}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	v, err := utils.GetPathValues(strings.Split(r.URL.Path, "/"), 0)
	if err != nil {
		utils.NewError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.q.GetUser(context.Background(), int32(v.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			utils.NewError(w, http.StatusNotFound, err.Error())
			return
		}

		utils.NewError(w, http.StatusInternalServerError, "error when getting user")
		return
	}

	utils.NewResponse(w, http.StatusOK, user)
	return
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.NewResponse(w, http.StatusMethodNotAllowed, consts.ResCreateUser{
			User:       consts.EmptyUser,
			FormErrors: consts.EmptyCreateUserErrors,
			Error:      consts.UnsupportedMethod.Error(),
		})
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))
	errs, valid := utils.ValidateNewUser(email, username, password)
	if !valid {
		utils.NewResponse(w, http.StatusUnprocessableEntity, consts.ResCreateUser{
			User: sqlc.User{
				UserID:   0,
				Email:    email,
				Username: username,
				Password: "",
				Admin:    false,
			},
			FormErrors: errs,
			Error:      "",
		})
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreateUser{
			User:       consts.EmptyUser,
			FormErrors: consts.EmptyCreateUserErrors,
			Error:      err.Error(),
		})
	}
	_, err = h.q.CreateUser(context.Background(), sqlc.CreateUserParams{
		Email:    email,
		Username: username,
		Password: string(hashedPass),
		Admin:    false,
	})
	if err != nil {
		utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreateUser{
			User:       consts.EmptyUser,
			FormErrors: consts.EmptyCreateUserErrors,
			Error:      err.Error(),
		})
		return
	}

	user, err := h.q.GetUserByEmail(context.Background(), r.FormValue("email"))
	if err != nil {
		utils.NewResponse(w, http.StatusInternalServerError, consts.ResCreateUser{
			User:       consts.EmptyUser,
			FormErrors: consts.EmptyCreateUserErrors,
			Error:      err.Error(),
		})
		return
	}

	utils.NewResponse(w, http.StatusCreated, consts.ResCreateUser{
		User:       user,
		FormErrors: consts.EmptyCreateUserErrors,
		Error:      "",
	})
	return
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if err := h.Authorizer(r, false); err != nil {
		utils.NewError(w, http.StatusUnauthorized, err.Error())
		return
	}

	v, err := utils.GetPathValues(strings.Split(r.URL.Path, "/"), 0)
	if err != nil {
		utils.NewError(w, http.StatusBadRequest, err.Error())
		return
	}

	c, err := r.Cookie("sid")
	if err != nil {
		utils.NewError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s, err := h.q.GetSession(context.Background(), c.Value)
	if err != nil {
		utils.NewError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if int(s.UserID) != v.Id {
		utils.NewError(w, http.StatusUnauthorized, "Permission not granted")
		return
	}

	exc, err := h.q.DeleteUser(context.Background(), int32(v.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			utils.NewError(w, http.StatusNotFound, err.Error())
			return
		}

		utils.NewError(w, http.StatusInternalServerError, "error when deleting user")
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
