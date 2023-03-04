package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jobutterfly/olives/utils"
)

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	v, err := utils.GetPathValues(strings.Split(r.URL.Path, "/"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`"msg" : "error while parsing path"`))
	}

	user, err := h.q.GetUser(context.Background(), int32(v.Id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`"msg" : "error when getting user"`))
	}

	jsonBytes, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`"msg" : "error when parsing user"`))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
