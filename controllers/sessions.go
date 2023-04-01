package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jobutterfly/olives/sqlc"
	"github.com/jobutterfly/olives/utils"
)

func (h *Handler) CreateSession(userId int32) (string, error) {
	sessionId, err := uuid.NewUUID()

	_, err = h.q.CreateSession(context.Background(), sqlc.CreateSessionParams{
		SessionID: sessionId.String(),
		UserID:    userId,
	})
	if err != nil {
		return "", err
	}

	return sessionId.String(), err
}

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("sid")
		if err != nil {
			if err == http.ErrNoCookie {
				utils.NewError(w, http.StatusBadRequest, err.Error())
			}
			utils.NewError(w, http.StatusInternalServerError, err.Error())
		}

		valid, err := h.VerifySession(c.Value)
		if err != nil {
			utils.NewError(w, http.StatusInternalServerError, err.Error())
		}
		if !valid {
			utils.NewError(w, http.StatusBadRequest, "invalid session")
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("sid")

		valid, err := h.VerifyAdmin(c.Value)
		if err != nil {
			utils.NewError(w, http.StatusInternalServerError, err.Error())
		}

		if !valid {
			utils.NewError(w, http.StatusBadRequest, "User not an admin")
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) ExtenderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("sid")
		if err != nil {
			if err == http.ErrNoCookie {
				next.ServeHTTP(w, r)
			}
			utils.NewError(w, http.StatusInternalServerError, err.Error())
		}

		valid, err := h.VerifySession(c.Value)
		if err != nil {
			utils.NewError(w, http.StatusInternalServerError, err.Error())
		}
		if !valid {
			next.ServeHTTP(w, r)
		}

		_, err = h.q.UpdateSession(context.Background(), c.Value)
		if err != nil {
			utils.NewError(w, http.StatusInternalServerError, err.Error())
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) VerifySession(sessiondId string) (bool, error) {
	session, err := h.q.GetSession(context.Background(), sessiondId)
	if err != nil {
		return false, err
	}

	now := time.Now()
	if session.LastAccess.Unix() < now.Add(time.Hour*-24).Unix() {
		_, err := h.q.DeleteSession(context.Background(), sessiondId)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func (h *Handler) VerifyAdmin(sessiondId string) (bool, error) {
	session, err := h.q.GetSession(context.Background(), sessiondId)
	if err != nil {
		return false, err
	}

	if session.Admin.Bool {
		return true, nil
	}

	return false, nil
}
