package utils

import (
	"database/sql"
	"math/rand"
	"strings"
	"time"

	"github.com/jobutterfly/olives/sqlc"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func randomString(n int) string {
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
		Email:    randomString(5) + "@" + randomString(6) + ".com",
		Username: randomString(10),
		Password: randomString(25),
	}
}

func RandomPost() sqlc.Post {
	rand.Seed(time.Now().Unix())
	return sqlc.Post{
		PostID: int32(rand.Intn(100)),
		Title: randomString(100),
		Text: randomString(1000),
		CreatedAt: time.Now(),
		UserID: int32(rand.Intn(100)),
		ImageID: sql.NullInt32{
			Int32: int32(rand.Intn(12)),
			Valid: true,
		},
		SuboliveID: int32(rand.Intn(5)),
	}
}

