package utils

import (
	"database/sql"
	"math/rand"
	"strings"
	"time"

	"github.com/jobutterfly/olives/sqlc"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	rand.Seed(time.Now().Unix())

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomUser() sqlc.User {
	rand.Seed(time.Now().Unix())
	return sqlc.User{
		UserID:   int32(rand.Intn(100)),
		Email:    RandomString(5) + "@" + RandomString(6) + ".com",
		Username: RandomString(10),
		Password: RandomString(25),
	}
}

func RandomPost() sqlc.Post {
	rand.Seed(time.Now().Unix())
	return sqlc.Post{
		PostID: int32(rand.Intn(100)),
		Title: RandomString(100),
		Text: RandomString(1000),
		CreatedAt: time.Now(),
		UserID: int32(rand.Intn(100)),
		ImageID: sql.NullInt32{
			Int32: int32(rand.Intn(12)),
			Valid: true,
		},
		SuboliveID: int32(rand.Intn(5)),
	}
}

