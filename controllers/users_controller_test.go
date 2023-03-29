package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/jobutterfly/olives/consts"
	"github.com/jobutterfly/olives/sqlc"
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

	testCases := []GetTestCase{
		{
			Name:         "successful get user request",
			Req:          httptest.NewRequest(http.MethodGet, "/users/"+strconv.Itoa(int(userId)), nil),
			ExpectedRes:  user,
			ExpectedCode: http.StatusOK,
		},
		{
			Name: "failed request for non existing user",
			Req:  httptest.NewRequest(http.MethodGet, "/users/"+strconv.Itoa(1000000), nil),
			ExpectedRes: consts.ErrorMessage{
				Msg: sql.ErrNoRows.Error(),
			},
			ExpectedCode: http.StatusNotFound,
		},
		{
			Name: "failed request for wrong path",
			Req:  httptest.NewRequest(http.MethodGet, "/users/banana", nil),
			ExpectedRes: consts.ErrorMessage{
				Msg: consts.PathNotAnInteger.Error(),
			},
			ExpectedCode: http.StatusBadRequest,
		},
	}

	TestGet(t, testCases, Th.GetUser)
}

func TestCreateUser(t *testing.T) {
	if err := Start(); err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	newestUser, err := Th.q.GetNewestUser(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	newUser := utils.RandomUser()
	newUser.UserID = newestUser.UserID + 1
	firstExpectedRes := consts.ResCreateUser{
		User:       newUser,
		FormErrors: consts.EmptyCreateUserErrors,
		Error:      "",
	}
	firstReq, err := NewPostRequestUser(newUser, "/users/create")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	secondUser := newUser
	secondUser.Email = "notanemail"
	secErrs, _ := utils.ValidateNewUser(secondUser.Email, secondUser.Username, secondUser.Password)
	secondExpectedRes := consts.ResCreateUser{
		User: sqlc.User{
			UserID:   0,
			Email:    secondUser.Email,
			Username: secondUser.Username,
			Password: "",
			Admin:    false,
		},
		FormErrors: secErrs,
		Error:      "",
	}
	secondReq, err := NewPostRequestUser(secondUser, "/users/create")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	thirdUser := newUser
	thirdUser.Username = "shor"
	thirdErrs, _ := utils.ValidateNewUser(thirdUser.Email, thirdUser.Username, thirdUser.Password)
	thirdExpectedRes := consts.ResCreateUser{
		User: sqlc.User{
			UserID:   0,
			Email:    thirdUser.Email,
			Username: thirdUser.Username,
			Password: "",
			Admin:    false,
		},
		FormErrors: thirdErrs,
		Error:      "",
	}
	thirdReq, err := NewPostRequestUser(thirdUser, "/users/create")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	fourthUser := newUser
	fourthUser.Password = "shor"
	fourthErrs, _ := utils.ValidateNewUser(fourthUser.Email, fourthUser.Username, fourthUser.Password)
	fourthExpectedRes := consts.ResCreateUser{
		User: sqlc.User{
			UserID:   0,
			Email:    fourthUser.Email,
			Username: fourthUser.Username,
			Password: "",
			Admin:    false,
		},
		FormErrors: fourthErrs,
		Error:      "",
	}
	fourthReq, err := NewPostRequestUser(fourthUser, "/users/create")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	testCases := []PostTestCase{
		{
			Name:         "successful create user request",
			Req:          firstReq,
			ExpectedRes:  firstExpectedRes,
			ExpectedCode: http.StatusCreated,
			TestAfter: AfterRes{
				Valid: true,
				Type:  "user",
			},
		},
		{
			Name:         "invalid email user request",
			Req:          secondReq,
			ExpectedRes:  secondExpectedRes,
			ExpectedCode: http.StatusUnprocessableEntity,
			TestAfter: AfterRes{
				Valid: false,
				Type:  "",
			},
		},
		{
			Name:         "invalid username user request",
			Req:          thirdReq,
			ExpectedRes:  thirdExpectedRes,
			ExpectedCode: http.StatusUnprocessableEntity,
			TestAfter: AfterRes{
				Valid: false,
				Type:  "",
			},
		},
		{
			Name:         "invalid password user request",
			Req:          fourthReq,
			ExpectedRes:  fourthExpectedRes,
			ExpectedCode: http.StatusUnprocessableEntity,
			TestAfter: AfterRes{
				Valid: false,
				Type:  "",
			},
		},
	}

	TestPost(t, testCases, Th.CreateUser)
}

func TestDeleteUser(t *testing.T) {
	if err := Start(); err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	newestUser, err := Th.q.GetNewestUser(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	firstReq := httptest.NewRequest(http.MethodDelete, "/users/delete/"+strconv.Itoa(int(newestUser.UserID)), nil)
	secondReq := httptest.NewRequest(http.MethodDelete, "/users/delete/"+strconv.Itoa(1000000), nil)

	testCases := []GetTestCase{
		{
			Name:         "successful delete of user",
			Req:          firstReq,
			ExpectedRes:  "",
			ExpectedCode: http.StatusOK,
		},
		{
			Name: "failed request for non existing user",
			Req:  secondReq,
			ExpectedRes: consts.ErrorMessage{
				Msg: sql.ErrNoRows.Error(),
			},
			ExpectedCode: http.StatusNotFound,
		},
	}

	TestGet(t, testCases, Th.DeleteUser)
}
