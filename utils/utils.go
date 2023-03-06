package utils

import (
	"errors"
	"strconv"
	"math/rand"
	"strings"

	"github.com/jobutterfly/olives/sqlc"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

type PathInfo struct {
    Id int
}

func GetPathValues(ps []string) (PathInfo, error){
    r := PathInfo{
	Id: 0,
    }

    if len(ps) > 3 {
	if ps[3] != "" {
	    err := errors.New("not found")
	    return r, err
	}
    }

    id, err := strconv.Atoi(ps[2])
    if err != nil {
	err := errors.New("not an integer")
	return r, err
    }
    r.Id = id

    return r, err
}


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
