package utils

import (
	"net/mail"
	"strconv"

	"github.com/jobutterfly/olives/consts"
)


type PathInfo struct {
	Id int
}

func GetPathValues(ps []string, offset int) (PathInfo, error) {
	r := PathInfo{
		Id: 0,
	}

	if len(ps) > 3 + offset {
		if ps[3 + offset] != "" {
			err := consts.PathNotFound
			return r, err
		}
	}

	id, err := strconv.Atoi(ps[2 + offset])
	if err != nil {
		err := consts.PathNotAnInteger
		return r, err
	}
	r.Id = id

	return r, err
}

func ValidateNewUser(email string, username string, password string) (errs [3]consts.FormInputError, valid bool) {
	valid = true
	errs = consts.EmptyCreateUserErrors

	_, err := mail.ParseAddress(email)
	if err != nil {
		errs[0].Bool = true
		errs[0].Message = "Invalid email address"
		valid = false
	}

	if username == "" {
		errs[1].Bool = true
		errs[1].Message = "This field is required"
		valid = false
	}

	if len(username) < 5 {
		errs[1].Bool = true
		errs[1].Message = "This field must be greater than 6 characters"
		valid = false
	}

	if password == "" {
		errs[2].Bool = true
		errs[2].Message = "This field is required"
		valid = false
	}

	if len(password) < 5 {
		errs[2].Bool = true
		errs[2].Message = "This field must be greater than 6 characters"
		valid = false
	}

	return errs, valid
}

