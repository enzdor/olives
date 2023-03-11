package consts

import (
	"errors"
)

type ErrorMessage struct {
	Msg string `json:"msg"`
}

var JsonParseError = `{"msg":"error while parsing"}`

var PathNotAnInteger = errors.New("path: not an integer")

var PathNotFound = errors.New("path: not found")

type FormInputError struct {
	Bool bool `json:"bool"`
	Message string `json:"msg"`
	Field string `json:"field"`
}

