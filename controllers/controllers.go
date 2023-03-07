package controllers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jobutterfly/olives/database"
	"github.com/joho/godotenv"
)

var Th *Handler

type GetTestCase struct {
	Name         string
	Req          *http.Request
	ExpectedRes  []byte
	ExpectedCode int
}

type PostTestCase struct {
	Name         string
	Req          *http.Request
	ExpectedRes  []byte
	ExpectedCode int
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

			if string(resBody) != string(tc.ExpectedRes) {
				t.Errorf("expected response body to be equal to: \n%v\ngot:\n%v", string(tc.ExpectedRes), string(resBody))
				return
			}
		})
	}

}
