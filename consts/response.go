package consts

import "github.com/jobutterfly/olives/sqlc"

type ResCreateUser struct {
	User sqlc.User `json:"user"`
	Errors [3]FormInputError `json:"errors"`
}

