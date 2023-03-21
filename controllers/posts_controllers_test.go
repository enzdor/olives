package controllers

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/jobutterfly/olives/consts"
	"github.com/jobutterfly/olives/sqlc"
	"github.com/jobutterfly/olives/utils"
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

	testCases := []GetTestCase{
		{
			Name:         "successful get post request",
			Req:          httptest.NewRequest(http.MethodGet, "/posts/"+strconv.Itoa(int(postId)), nil),
			ExpectedRes:  post,
			ExpectedCode: http.StatusOK,
		},
		{
			Name: "failed request for non existing post",
			Req:  httptest.NewRequest(http.MethodGet, "/posts/"+strconv.Itoa(1000000), nil),
			ExpectedRes: consts.ErrorMessage{
				Msg: sql.ErrNoRows.Error(),
			},
			ExpectedCode: http.StatusNotFound,
		},
		{
			Name: "failed request for wrong path",
			Req:  httptest.NewRequest(http.MethodGet, "/posts/banana", nil),
			ExpectedRes: consts.ErrorMessage{
				Msg: consts.PathNotAnInteger.Error(),
			},
			ExpectedCode: http.StatusBadRequest,
		},
	}

	TestGet(t, testCases, Th.GetPost)
}

func TestGetSubolivePosts(t *testing.T) {
	suboliveId := int32(2)
	firstPosts, err := Th.q.GetSubolivePosts(context.Background(), sqlc.GetSubolivePostsParams{
		Offset:     0,
		SuboliveID: suboliveId,
	})
	if err != nil {
		t.Errorf("expected no errors, got %v", err)
		return
	}

	secondPosts, err := Th.q.GetSubolivePosts(context.Background(), sqlc.GetSubolivePostsParams{
		Offset:     10,
		SuboliveID: suboliveId,
	})
	if err != nil {
		t.Errorf("expected no errors, got %v", err)
		return
	}

	testCases := []GetTestCase{
		{
			Name:         "successful get posts request",
			Req:          httptest.NewRequest(http.MethodGet, "/posts/subolive/"+strconv.Itoa(int(suboliveId))+"?page=0", nil),
			ExpectedRes:  firstPosts,
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "successful get posts request second page",
			Req:          httptest.NewRequest(http.MethodGet, "/posts/subolive/"+strconv.Itoa(int(suboliveId))+"?page=1", nil),
			ExpectedRes:  secondPosts,
			ExpectedCode: http.StatusOK,
		},
		{
			Name: "unsuccessful get posts request subolive id not number",
			Req:  httptest.NewRequest(http.MethodGet, "/posts/subolive/banana?page=1", nil),
			ExpectedRes: consts.ErrorMessage{
				Msg: consts.PathNotAnInteger.Error(),
			},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name: "unsuccessful get posts request subolive id not exist",
			Req:  httptest.NewRequest(http.MethodGet, "/posts/subolive/"+strconv.Itoa(int(10000))+"?page=1", nil),
			ExpectedRes: consts.ErrorMessage{
				Msg: consts.SuboliveNonExistant.Error(),
			},
			ExpectedCode: http.StatusNotFound,
		},
		{
			Name: "unsuccessful get posts request not existant page",
			Req:  httptest.NewRequest(http.MethodGet, "/posts/subolive/"+strconv.Itoa(int(suboliveId))+"?page=1000", nil),
			ExpectedRes: consts.ErrorMessage{
				Msg: sql.ErrNoRows.Error(),
			},
			ExpectedCode: http.StatusNotFound,
		},
		{
			Name: "unsuccessful get posts request page not int",
			Req:  httptest.NewRequest(http.MethodGet, "/posts/subolive/"+strconv.Itoa(int(suboliveId))+"?page=banana", nil),
			ExpectedRes: consts.ErrorMessage{
				Msg: consts.PageNotAnInteger.Error(),
			},
			ExpectedCode: http.StatusInternalServerError,
		},
	}

	TestGet(t, testCases, Th.GetSubolivePosts)
}

func TestCreatePost(t *testing.T) {
	newestPost, err := Th.q.GetNewestPost(context.Background())
	if err != nil {
		t.Errorf("expected no errors, got %v", err)
		return
	}

	newPost := utils.RandomPost()
	newPost.PostID = newestPost.PostID + 1

	firstExpectedRes := consts.ResCreatedPost{
		Post:   newPost,
		Errors: consts.EmptyCreatePostErrors,
	}

	newPost2 := sqlc.Post{
		PostID:     0,
		Title:      "",
		Text:       "",
		CreatedAt:  time.Now(),
		UserID:     20,
		SuboliveID: 4,
		ImageID: sql.NullInt32{
			Int32: 0,
			Valid: false,
		},
	}
	errs2 := consts.EmptyCreatePostErrors
	errs2[0].Bool = true
	errs2[0].Message = "This field must be greater than 6 characters"
	errs2[1].Bool = true
	errs2[1].Message = "This field must be greater than 6 characters"
	errs2[2].Bool = true
	errs2[2].Message = "File size greater than 512 kilobytes. Choose a smaller file."
	secondExpectedRes := consts.ResCreatedPost{
		Post:   newPost2,
		Errors: errs2,
	}

	pr, pw := io.Pipe()
	form := multipart.NewWriter(pw)
	go NewPost(t, pw, form, "test.png", newPost)

	firstReq, err := http.NewRequest(http.MethodPost, "/posts", pr)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	firstReq.Header.Set("Content-Type", form.FormDataContentType())

	pr2, pw2 := io.Pipe()
	form2 := multipart.NewWriter(pw2)
	go NewPost(t, pw2, form2, "cheese.png", newPost2)

	secondReq, err := http.NewRequest(http.MethodPost, "/posts", pr2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	secondReq.Header.Set("Content-Type", form2.FormDataContentType())

	thirdPost := utils.RandomPost()
	thirdPost.PostID = newestPost.PostID + 2
	thirdBody := bytes.NewReader([]byte("title=" + thirdPost.Title + "&text=" + thirdPost.Text + "&user_id=" + strconv.Itoa(int(thirdPost.UserID)) + "&subolive_id=" + strconv.Itoa(int(thirdPost.SuboliveID)) + "&image="))
	thirdReq := httptest.NewRequest(http.MethodPost, "/posts", thirdBody)
	thirdReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	thirdExpectedRes := consts.ResCreatedPost{
		Post:   thirdPost,
		Errors: consts.EmptyCreatePostErrors,
	}

	fourthPost := newPost2
	fourthPost.Text = "bla"
	fourthPost.Title = "a valid title"
	fourthErrs := consts.EmptyCreatePostErrors
	fourthErrs[1].Bool = true
	fourthErrs[1].Message = "This field must be greater than 6 characters"
	fourthBody := bytes.NewReader([]byte("title=" + fourthPost.Title + "&text=" + fourthPost.Text + "&user_id=" + strconv.Itoa(int(fourthPost.UserID)) + "&subolive_id=" + strconv.Itoa(int(fourthPost.SuboliveID)) + "&image="))
	fourthReq := httptest.NewRequest(http.MethodPost, "/posts", fourthBody)
	fourthReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	fourthExpectedRes := consts.ResCreatedPost{
		Post:   fourthPost,
		Errors: fourthErrs,
	}

	fifthPost := newPost2
	fifthPost.Text = "a valid text bla bla"
	fifthPost.Title = utils.RandomString(300)
	fifthErrs := consts.EmptyCreatePostErrors
	fifthErrs[0].Bool = true
	fifthErrs[0].Message = "This field must have less than 255 characters"
	fifthBody := bytes.NewReader([]byte("title=" + fifthPost.Title + "&text=" + fifthPost.Text + "&user_id=" + strconv.Itoa(int(fifthPost.UserID)) + "&subolive_id=" + strconv.Itoa(int(fifthPost.SuboliveID)) + "&image="))
	fifthReq := httptest.NewRequest(http.MethodPost, "/posts", fifthBody)
	fifthReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	fifthExpectedRes := consts.ResCreatedPost{
		Post:   fifthPost,
		Errors: fifthErrs,
	}

	sixthPost := newPost
	sixthPost.PostID = 0
	sixthPost.Title = "sh"
	sixthPost.Text = ""
	sixthPost.ImageID = sql.NullInt32{
		Int32: 0,
		Valid: false,
	}
	sixthErrs := [3]consts.FormInputError{
		{
			Bool: true,
			Message: "This field must be greater than 6 characters",
			Field: "title",
		},
		{
			Bool: true,
			Message: "This field is required",
			Field: "text",
		},
		{
			Bool: false,
			Message: "",
			Field: "image",
		},
	}
	pr6, pw6 := io.Pipe()
	form6 := multipart.NewWriter(pw6)
	go NewPost(t, pw6, form6, "test.png", sixthPost)

	sixthReq, err := http.NewRequest(http.MethodPost, "/posts", pr6)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	sixthReq.Header.Set("Content-Type", form6.FormDataContentType())
	sixthExpectedRes := consts.ResCreatedPost{
		Post: sixthPost,
		Errors: sixthErrs,
	}

	testCases := []PostTestCase{
		{
			Name:         "successful post post",
			Req:          firstReq,
			ExpectedRes:  firstExpectedRes,
			ExpectedCode: http.StatusCreated,
			TestAfter: AfterRes{
				Valid: true,
				Type:  "post",
			},
		},
		{
			Name:         "unsuccessful post post, image too large",
			Req:          secondReq,
			ExpectedRes:  secondExpectedRes,
			ExpectedCode: http.StatusUnprocessableEntity,
			TestAfter: AfterRes{
				Valid: false,
				Type:  "",
			},
		},
		{
			Name:         "succesful post without image",
			Req:          thirdReq,
			ExpectedRes:  thirdExpectedRes,
			ExpectedCode: http.StatusCreated,
			TestAfter: AfterRes{
				Valid: true,
				Type:  "post",
			},
		},
		{
			Name:         "unsuccesful post without image, text too short",
			Req:          fourthReq,
			ExpectedRes:  fourthExpectedRes,
			ExpectedCode: http.StatusUnprocessableEntity,
			TestAfter: AfterRes{
				Valid: false,
				Type:  "",
			},
		},
		{
			Name:         "unsuccesful post without image, title too long",
			Req:          fifthReq,
			ExpectedRes:  fifthExpectedRes,
			ExpectedCode: http.StatusUnprocessableEntity,
			TestAfter: AfterRes{
				Valid: false,
				Type:  "",
			},
		},
		{
			Name:         "unsuccessful post post with image, title short text required",
			Req:          sixthReq,
			ExpectedRes:  sixthExpectedRes,
			ExpectedCode: http.StatusUnprocessableEntity,
			TestAfter: AfterRes{
				Valid: false,
				Type:  "",
			},
		},
	}

	TestPost(t, testCases, Th.CreatePost)
}
