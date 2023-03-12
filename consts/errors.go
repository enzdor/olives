package consts

import (
	"errors"
)

type ErrorMessage struct {
	Msg string `json:"msg"`
}

var JsonParseError = `{"msg":"error while parsing"}`

var PathNotAnInteger = errors.New("path: id not an integer")

var PathNotFound = errors.New("path: not found")

var PageNotAnInteger = errors.New("path: page not an integer")

var SuboliveNonExistant = errors.New("path: subolive does'nt exist")

type FormInputError struct {
	Bool bool `json:"bool"`
	Message string `json:"msg"`
	Field string `json:"field"`
}

var EmptyCreateUserErrors = [3]FormInputError {
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
}
