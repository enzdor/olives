package consts

import (
	"errors"
)

type ErrorMessage struct {
	Msg string `json:"msg"`
}

type FormInputError struct {
	Bool    bool   `json:"bool"`
	Message string `json:"msg"`
	Field   string `json:"field"`
}

var (
	JsonParseError        = `{"msg":"error while parsing"}`
	PathNotAnInteger      = errors.New("path: id not an integer")
	PathNotFound          = errors.New("path: not found")
	PageNotAnInteger      = errors.New("path: page not an integer")
	SuboliveNonExistant   = errors.New("path: subolive does'nt exist")
	UnsupportedMethod     = errors.New("Unsupported method")
	EmptyCreateUserErrors = [3]FormInputError{
		{
			Bool:    false,
			Message: "",
			Field:   "email",
		},
		{
			Bool:    false,
			Message: "",
			Field:   "username",
		},
		{
			Bool:    false,
			Message: "",
			Field:   "password",
		},
	}
	EmptyCreatePostErrors = [3]FormInputError{
		{
			Bool:    false,
			Message: "",
			Field:   "title",
		},
		{
			Bool:    false,
			Message: "",
			Field:   "text",
		},
		{
			Bool:    false,
			Message: "",
			Field:   "image",
		},
	}
	EmptyCreateCommentErrors = [2]FormInputError{
		{
			Bool:    false,
			Message: "",
			Field:   "text",
		},
		{
			Bool:    false,
			Message: "",
			Field:   "image",
		},
	}
)
