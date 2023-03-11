package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/jobutterfly/olives/consts"
	"github.com/jobutterfly/olives/sqlc"
	"github.com/jobutterfly/olives/utils"
)

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	v, err := utils.GetPathValues(strings.Split(r.URL.Path, "/"))
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

/*
	TODO: Add validation for form values
*/

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))
	errs, valid := utils.ValidateNewUser(email, username, password)
	if !valid {
		utils.NewErrorBody(w, http.StatusUnprocessableEntity, consts.ResCreateUser{
			User: sqlc.User {
				UserID: 0,
				Email: email,
				Username: username,
				Password: "",
			},
			Errors: errs,
		})
		return
	}
	_, err := h.q.CreateUser(context.Background(), sqlc.CreateUserParams{
		Email: email,
		Username: username,
		Password: password,
	})
	if err != nil {
		utils.NewError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.q.GetUserByEmail(context.Background(), r.FormValue("email"))
	if err != nil {
		utils.NewError(w, http.StatusInternalServerError, err.Error())
	}

	utils.NewResponse(w, http.StatusCreated, consts.ResCreateUser{
		User: user,
		Errors: [3]consts.FormInputError {
			{
				Bool: false,
				Message: "",
				Field: "email",
			},
			{
				Bool: false,
				Message: "",
				Field: "username",
			},
			{
				Bool: false,
				Message: "",
				Field: "password",
			},
		},
	})
	return
}















