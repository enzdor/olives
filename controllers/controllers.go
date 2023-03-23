package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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
	ExpectedRes  any
	ExpectedCode int
}

type AfterRes struct {
	Valid bool
	Type  string
}

type PostTestCase struct {
	Name         string
	Req          *http.Request
	ExpectedRes  any
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

			expectedResJson, err := json.Marshal(tc.ExpectedRes)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
				return
			}

			if string(resBody) != string(expectedResJson) {
				t.Errorf("expected response body to be equal to: \n%v\ngot:\n%v", string(expectedResJson), string(resBody))
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
						Post:       resPost,
						FormErrors: consts.EmptyCreatePostErrors,
					}

					resJson, err := json.Marshal(res)
					if err != nil {
						t.Errorf("expected no errors, got %v", err)
						return
					}

					if string(resBody) != string(resJson) {
						t.Errorf("expected response body to be equal to: \n%v\ngot:\n%v", string(resJson), string(resBody))
						return
					}
				}
			} else {
				expectedResJson, err := json.Marshal(tc.ExpectedRes)
				if err != nil {
					t.Errorf("expected no errors, got %v", err)
					return
				}

				if string(resBody) != string(expectedResJson) {
					t.Errorf("expected response body to be equal to: \n%v\ngot:\n%v", string(expectedResJson), string(resBody))
					return
				}
			}
		})
	}
}

func NewPostRequestPostImage(t *testing.T, writer *io.PipeWriter, form *multipart.Writer, path string, post sqlc.Post) {
	defer writer.Close()

	if err := form.WriteField("title", post.Title); err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	if err := form.WriteField("text", post.Text); err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	if err := form.WriteField("user_id", strconv.Itoa(int(post.UserID))); err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	if err := form.WriteField("subolive_id", strconv.Itoa(int(post.SuboliveID))); err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	file, err := os.Open("../" + path)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	defer file.Close()

	w, err := form.CreateFormFile("image", path)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	_, err = io.Copy(w, file)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	form.Close()
}

func NewPostRequestCommentImage(t *testing.T, writer *io.PipeWriter, form *multipart.Writer, path string, comment sqlc.Comment) {
	defer writer.Close()

	if err := form.WriteField("text", comment.Text); err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	if err := form.WriteField("user_id", strconv.Itoa(int(comment.UserID))); err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	if err := form.WriteField("post_id", strconv.Itoa(int(comment.PostID))); err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	file, err := os.Open("../" + path)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	defer file.Close()

	w, err := form.CreateFormFile("image", path)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	_, err = io.Copy(w, file)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	form.Close()
}

func NewPostRequestUser(u sqlc.User, target string) (*http.Request, error) {
	body := bytes.NewReader([]byte("username=" + u.Username + "&email=" + u.Email + "&password=" + u.Password))
	req, err := http.NewRequest(http.MethodPost, target, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func NewPostRequestPost(p sqlc.Post, target string) (*http.Request, error) {
	body := bytes.NewReader([]byte("title=" + p.Title + "&text=" + p.Text + "&user_id=" + strconv.Itoa(int(p.UserID)) + "&subolive_id=" + strconv.Itoa(int(p.SuboliveID)) + "&image="))
	req, err := http.NewRequest(http.MethodPost, target, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func NewPostRequestComment(c sqlc.Comment, target string) (*http.Request, error) {
	body := bytes.NewReader([]byte("text=" + c.Text + "&user_id=" + strconv.Itoa(int(c.UserID)) + "&post_id=" + strconv.Itoa(int(c.PostID)) + "&image="))
	req, err := http.NewRequest(http.MethodPost, target, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}
