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
	"github.com/jobutterfly/olives/utils"
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

/*
	FIXME: instead of checking with createdUser, query the testdb for the 
	latest created user use it to check if res is correct. the problem is
	that the id of the created user is different from the one actually 
	created
*/

func TestCreateUser(t *testing.T) {
	newestUser, err := Th.q.GetNewestUser(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	newUser := utils.RandomUser()
	newUser.UserID = newestUser.UserID + 1

	expectedRes := consts.ResCreateUser {
		User: newUser,
		Errors: [3]consts.FormInputError{
			{
				Bool: false,
				Message: "",
				Field: "email",
			},
			{
				Bool: false,
				Message: "",
				Field: "username",

			},
			{
				Bool: false,
				Message: "",
				Field: "password",
			},
		},
	}
	jsonRes, err := json.Marshal(expectedRes)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	firstBody := bytes.NewReader([]byte("username=" + newUser.Username + "&email=" + newUser.Email + "&password=" + newUser.Password))
	firstReq := httptest.NewRequest(http.MethodPost, "/users", firstBody)
	firstReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	testCases := []PostTestCase{
		{
			Name:         "successful create user request",
			Req:          firstReq,
			ExpectedRes:  jsonRes,
			ExpectedCode: http.StatusCreated,
		},
	}

	TestPost(t, testCases, Th.CreateUser)
}
