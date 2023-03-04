package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/jobutterfly/olives/database"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)


var th *Handler

func start() error{
    if err := godotenv.Load("../.env"); err != nil {
	return err
    }
    user := os.Getenv("DBUSER")
    pass := os.Getenv("DBPASS")
    name := os.Getenv("TESTDBNAME")
    key := os.Getenv("JWTKEY")

    db := database.NewDB(user, pass, name)

    th = NewHandler(db, key)

    return nil
}

func TestGetUser(t *testing.T) {
	if err := start(); err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	userId := int32(9)
	user, err := th.q.GetUser(context.Background(), userId)
	if err != nil {
		t.Errorf("expected no errors, got %v", err)
		return
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	testCases := []struct {
		req *http.Request
		expected []byte
	} {
		{
			req: httptest.NewRequest(http.MethodGet, "/users/" + strconv.Itoa(int(userId)), nil),
			expected: jsonUser,
		},
	}

	for _, i := range testCases {
		w := httptest.NewRecorder()

		th.GetUser(w, i.req)
		res := w.Result()
		defer res.Body.Close()

		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
			return
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status to be %v, got %v", http.StatusOK, res.StatusCode)
			return
		}

		if string(resBody) != string(i.expected) {
			t.Errorf("expected response body to be equal to: \n%v\ngot:\n%v", string(i.expected), string(resBody))
			return
		}
	}
}











