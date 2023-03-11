package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/jobutterfly/olives/consts"
)

func TestGetPost(t *testing.T) {
	if err := Start(); err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	postId := int32(9)
	post, err := Th.q.GetPost(context.Background(), postId)
	if err != nil {
		t.Errorf("expected no errors, got %v", err)
		return
	}

	jsonPost, err := json.Marshal(post)
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
			Name:         "successful get post request",
			Req:          httptest.NewRequest(http.MethodGet, "/posts/"+strconv.Itoa(int(postId)), nil),
			ExpectedRes:  jsonPost,
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "failed request for non existing post",
			Req:          httptest.NewRequest(http.MethodGet, "/posts/"+strconv.Itoa(1000000), nil),
			ExpectedRes:  noRowMsg,
			ExpectedCode: http.StatusNotFound,
		},
		{
			Name:         "failed request for wrong path",
			Req:          httptest.NewRequest(http.MethodGet, "/posts/banana", nil),
			ExpectedRes:  notAnInt,
			ExpectedCode: http.StatusBadRequest,
		},
	}

	TestGet(t, testCases, Th.GetPost)
}
