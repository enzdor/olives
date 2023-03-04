package controllers

import (
	"database/sql"

	"github.com/jobutterfly/olives/sqlc"
)

type Handler struct {
	q  *sqlc.Queries
	key string
}

func NewHandler(db *sql.DB, key string) *Handler {
	queries := sqlc.New(db)

	return &Handler{
		q: queries,
		key: key,
	}
}
