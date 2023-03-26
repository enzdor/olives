package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jobutterfly/olives/sqlc"
)

func (h *Handler) CreateSession(userId int32) (string, error) {
	sessionId, err := uuid.NewUUID()

	_, err = h.q.CreateSession(context.Background(), sqlc.CreateSessionParams{
		SessionID: sessionId.String(),
		UserID: userId,
	})
	if err != nil {
		return "", err
	}

	return sessionId.String(), err
}

func (h *Handler) Authorizer(r *http.Request, needsAdmin bool) error {
	c, err := r.Cookie("sid")
	if err != nil {
		if err == http.ErrNoCookie {
			return errors.New("Cookie not found")
		}
		return err
	}

	if needsAdmin {
		isAdmin, err := h.VerifyAdmin(c.Value)
		if err != nil {
			return err
		}
		if !isAdmin {
			return errors.New("User is not admin")
		}
	}

	valid, err := h.VerifySession(c.Value)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("Session expired")
	}

	return nil
}

func (h *Handler) VerifySession(sessiondId string) (bool, error) {
	session, err := h.q.GetSession(context.Background(), sessiondId)
	if err != nil {
		return false, err
	}

	now := time.Now()
	if session.LastAccess.Unix() < now.Add(time.Hour * -24).Unix() {
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

func (h *Handler) ExtendSessionIfExists(r *http.Request) error {
	c, err := r.Cookie("sid")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil
		}
		return err
	}

	valid, err := h.VerifySession(c.Value)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("Session expired")
	}

	_, err = h.q.UpdateSession(context.Background(), c.Value)
	if err != nil {
		return err
	}

	return nil
}











