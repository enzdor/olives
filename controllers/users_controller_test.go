package controllers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/jobutterfly/olives/consts"
	"github.com/jobutterfly/olives/sqlc"
)

func TestGetUser(t *testing.T) {
	if err := Start(); err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	userId := int32(9)
	user, err := Th.q.GetUser(context.Background(), userId)
	if err != nil {
		t.Errorf("expected no errors, got %v", err)
		return
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	noRowMsg, err := json.Marshal(consts.ErrorMessage{
		Msg: sql.ErrNoRows.Error(),
	})
	if err != nil {
		t.Errorf("expected no errors, got %v", err)
		return
	}

	notAnInt, err := json.Marshal(consts.ErrorMessage{
		Msg: consts.PathNotAnInteger.Error(),
	})
	if err != nil {
		t.Errorf("expected no errors, got %v", err)
		return
	}

	testCases := []GetTestCase{
		{
			Name:         "successful get user request",
			Req:          httptest.NewRequest(http.MethodGet, "/users/"+strconv.Itoa(int(userId)), nil),
			ExpectedRes:  jsonUser,
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "failed request for non existing user",
			Req:          httptest.NewRequest(http.MethodGet, "/users/"+strconv.Itoa(1000000), nil),
			ExpectedRes:  noRowMsg,
			ExpectedCode: http.StatusNotFound,
		},
		{
			Name:         "failed request for wrong path",
			Req:          httptest.NewRequest(http.MethodGet, "/users/banana", nil),
			ExpectedRes:  notAnInt,
			ExpectedCode: http.StatusBadRequest,
		},
	}

	TestGet(t, testCases, Th.GetUser)
}

func TestCreateUser(t *testing.T) {

	createdUser := sqlc.User{
		UserID:   int32(101),
		Username: "banana",
		Email:    "banana@tree.com",
		Password: "supersecret",
	}

	jsonUser, err := json.Marshal(createdUser)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	firstBody := bytes.NewReader([]byte("username=" + createdUser.Username + "&email=" + createdUser.Email + "&password=" + createdUser.Password))
	firstReq := httptest.NewRequest(http.MethodPost, "/users", firstBody)
	firstReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	testCases := []PostTestCase{
		{
			Name:         "successful create user request",
			Req:          firstReq,
			ExpectedRes:  jsonUser,
			ExpectedCode: http.StatusCreated,
		},
	}

	TestPost(t, testCases, Th.CreateUser)
}
