package utils

import (
	"encoding/json"
	"net/http"

	"github.com/jobutterfly/olives/consts"
)

func NewError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	jsonBytes, err := json.Marshal(consts.ErrorMessage{Msg: msg})
	if err != nil {
		w.Write([]byte(consts.JsonParseError))
		return
	}
	w.Write(jsonBytes)
	return
}

func NewResponse(w http.ResponseWriter, status int, body any) {
	w.WriteHeader(status)
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(consts.JsonParseError))
		return
	}
	w.Write(jsonBytes)
	return
}

