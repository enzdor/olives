package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
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
	
	notFoundMsg, err := json.Marshal(struct{
		Msg string `json:"msg"`
	}{
		Msg: "user not found",
	})
	if err != nil {
		t.Errorf("expected no errors, got %v", err)
		return
	}


	testCases := []GetTestCase {
		{
			Name: "successful get user request",
			Req: httptest.NewRequest(http.MethodGet, "/users/" + strconv.Itoa(int(userId)), nil),
			ExpectedRes: jsonUser,
			ExpectedCode: http.StatusOK,
		}, 
		{
			Name: "failed request for non existing user",
			Req: httptest.NewRequest(http.MethodGet, "/users/" + strconv.Itoa(1000000), nil),
			ExpectedRes: notFoundMsg,
			ExpectedCode: http.StatusNotFound,
		},
	}

	TestGet(t, testCases, Th.GetUser)
}











