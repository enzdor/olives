package utils

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/jobutterfly/olives/consts"
	"github.com/jobutterfly/olives/sqlc"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

type PathInfo struct {
	Id int
}

func GetPathValues(ps []string) (PathInfo, error) {
	r := PathInfo{
		Id: 0,
	}

	if len(ps) > 3 {
		if ps[3] != "" {
			err := consts.PathNotFound
			return r, err
		}
	}

	id, err := strconv.Atoi(ps[2])
	if err != nil {
		err := consts.PathNotAnInteger
		return r, err
	}
	r.Id = id

	return r, err
}

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

func NewErrorBody(w http.ResponseWriter, status int, body any) {
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

// TODO: return object of errors for frontend

func ValidateNewUser(email string, username string, password string) (errs [3]consts.FormInputError, valid bool) {
	valid = true
	errs = consts.EmptyCreateUserErrors

	_, err := mail.ParseAddress(email)
	if err != nil {
		errs[0].Bool = true
		errs[0].Message = "Invalid email address"
		valid = false
	}

	if username == "" {
		errs[1].Bool = true
		errs[1].Message = "This field is required"
		valid = false
	}

	if len(username) < 5 {
		errs[1].Bool = true
		errs[1].Message = "This field must be greater than 6 characters"
		valid = false
	}

	if password == "" {
		errs[2].Bool = true
		errs[2].Message = "This field is required"
		valid = false
	}

	if len(password) < 5 {
		errs[2].Bool = true
		errs[2].Message = "This field must be greater than 6 characters"
		valid = false
	}

	return errs, valid
}





