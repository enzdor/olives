package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jobutterfly/olives/consts"
	"github.com/jobutterfly/olives/database"
	"github.com/jobutterfly/olives/sqlc"
	"github.com/jobutterfly/olives/utils"
	"github.com/joho/godotenv"
)

var Th *Handler

type GetTestCase struct {
	Name         string
	Req          *http.Request
	ExpectedRes  []byte
	ExpectedCode int
}

type AfterRes struct {
	Valid bool
	Type  string
}

type PostTestCase struct {
	Name         string
	Req          *http.Request
	ExpectedRes  []byte
	ExpectedCode int
	TestAfter    AfterRes
}

func Start() error {
	if err := godotenv.Load("../.env"); err != nil {
		return err
	}
	user := os.Getenv("DBUSER")
	pass := os.Getenv("DBPASS")
	name := os.Getenv("TESTDBNAME")
	key := os.Getenv("JWTKEY")

	db := database.NewDB(user, pass, name)

	Th = NewHandler(db, key)

	return nil
}

func TestGet(t *testing.T, testCases []GetTestCase, controller func(w http.ResponseWriter, r *http.Request)) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			w := httptest.NewRecorder()

			controller(w, tc.Req)
			res := w.Result()
			defer res.Body.Close()

			resBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
				return
			}

			if res.StatusCode != tc.ExpectedCode {
				t.Errorf("expected status to be %v, got %v", tc.ExpectedCode, res.StatusCode)
				return
			}

			if string(resBody) != string(tc.ExpectedRes) {
				t.Errorf("expected response body to be equal to: \n%v\ngot:\n%v", string(tc.ExpectedRes), string(resBody))
				return
			}
		})
	}
}

func TestPost(t *testing.T, testCases []PostTestCase, controller func(w http.ResponseWriter, r *http.Request)) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			w := httptest.NewRecorder()

			controller(w, tc.Req)
			res := w.Result()
			defer res.Body.Close()

			resBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
				return
			}

			if res.StatusCode != tc.ExpectedCode {
				t.Errorf("expected status to be %v, got %v", tc.ExpectedCode, res.StatusCode)
				return
			}

			if tc.TestAfter.Valid {
				switch tc.TestAfter.Type {
				case "post":
					post, err := Th.q.GetNewestPost(context.Background())
					if err != nil {
						utils.NewError(w, http.StatusInternalServerError, err.Error())
						return
					}

					resPost := sqlc.Post{
						PostID:     post.PostID,
						Title:      post.Title,
						Text:       post.Text,
						CreatedAt:  post.CreatedAt,
						UserID:     post.UserID,
						SuboliveID: post.SuboliveID,
						ImageID:    post.ImageID,
					}

					res := consts.ResCreatedPost{
						Post:   resPost,
						Errors: consts.EmptyCreatePostErrors,
					}

					resJson, err := json.Marshal(res)
					if err != nil {
						t.Errorf("expected no errors, got %v", err)
						return
					}

					if string(resBody) != string(resJson) {
						t.Errorf("banana expected response body to be equal to: \n%v\ngot:\n%v", string(resJson), string(resBody))
						return
					}
				}
			} else {
				if string(resBody) != string(tc.ExpectedRes) {
					t.Errorf("expected response body to be equal to: \n%v\ngot:\n%v", string(tc.ExpectedRes), string(resBody))
					return
				}
			}
		})
	}
}
