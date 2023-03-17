package consts

import "github.com/jobutterfly/olives/sqlc"

type ResCreateUser struct {
	User sqlc.User `json:"user"`
	Errors [3]FormInputError `json:"errors"`
}

type ResCreatedPost struct {
	Post sqlc.Post `json:"post"`
	Errors [3]FormInputError `json:"errors"`
}

const ITEMS_PER_PAGE = 10
