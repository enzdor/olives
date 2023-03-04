package utils

import (
	"strings"
	"math/rand"
	"github.com/jobutterfly/olives/sqlc"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomUser() sqlc.User {
	return sqlc.User{
		UserID: int32(rand.Intn(100)),
		Email: randomString(25),
		Username: randomString(10),
		Password: randomString(25),
	}
}
