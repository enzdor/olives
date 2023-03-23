package consts

import "github.com/jobutterfly/olives/sqlc"

type ResCreateUser struct {
	User       sqlc.User         `json:"user"`
	FormErrors [3]FormInputError `json:"form_errors"`
	Error      string            `json:"error"`
}

type ResCreatedPost struct {
	Post       sqlc.Post         `json:"post"`
	FormErrors [3]FormInputError `json:"form_errors"`
	Error      string            `json:"error"`
}

type ResCreatedComment struct {
	Comment    sqlc.Comment      `json:"comment"`
	FormErrors [2]FormInputError `json:"form_errors"`
	Error      string            `json:"error"`
}

const ITEMS_PER_PAGE = 10
