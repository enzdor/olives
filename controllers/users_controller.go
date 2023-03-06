package controllers

import (
	"context"
	"database/sql"
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
		return
	}

	user, err := h.q.GetUser(context.Background(), int32(v.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"msg":"user not found"}`))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`"msg" : "error when getting user"`))
		return
	}

	jsonBytes, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`"msg" : "error when parsing user"`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
