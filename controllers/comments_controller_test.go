package controllers

import (
	"context"
	"database/sql"
	"io"
	"log"
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

func TestCreateComment(t *testing.T) {
	newestComment, err := Th.q.GetNewestComment(context.Background())
	if err != nil {
		t.Errorf("expected no errors, got %v", err)
		return
	}
	log.Fatal("banana")

	newComment := utils.RandomComment()
	newComment.CommentID = newestComment.CommentID + 1

	firstExpectedRes := consts.ResCreatedComment{
		Post:       newComment,
		FormErrors: consts.EmptyCreateCommentErrors,
		Error:      "",
	}
	newComment2 := sqlc.Comment{
		CommentID: 0,
		Text:      "",
		CreatedAt: time.Now(),
		UserID:    20,
		ImageID: sql.NullInt32{
			Int32: 0,
			Valid: false,
		},
		PostID: 0,
	}
	secondExpectedRes := consts.ResCreatedComment{
		Post: newComment2,
		FormErrors: [2]consts.FormInputError{
			{
				Bool:    true,
				Message: "This field is required",
				Field:   "text",
			},
			{
				Bool:    true,
				Message: "File size greater than 512 kilobytes. Choose a smaller file.",
				Field:   "image",
			},
		},
		Error: "",
	}

	pr, pw := io.Pipe()
	form := multipart.NewWriter(pw)
	go NewPostRequestCommentImage(t, pw, form, "test.png", newComment)
	firstReq, err := http.NewRequest(http.MethodPost, "/comments", pr)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	firstReq.Header.Set("Content-Type", form.FormDataContentType())

	pr2, pw2 := io.Pipe()
	form2 := multipart.NewWriter(pw2)
	go NewPostRequestCommentImage(t, pw2, form2, "cheese.png", newComment2)
	secondReq, err := http.NewRequest(http.MethodPost, "/comments", pr2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	secondReq.Header.Set("Content-Type", form2.FormDataContentType())

	thirdComment := utils.RandomComment()
	thirdReq, err := NewPostRequestComment(thirdComment, "/comments")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	thirdExpectedRes := consts.ResCreatedComment{
		Post:       thirdComment,
		FormErrors: consts.EmptyCreateCommentErrors,
		Error:      "",
	}

	fourthComment := newComment2
	fourthComment.Text = "bla"
	fourthErrs := consts.EmptyCreateCommentErrors
	fourthErrs[1].Bool = true
	fourthErrs[1].Message = "This field must be greater than 6 characters"
	fourthReq, err := NewPostRequestComment(fourthComment, "/posts")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	fourthExpectedRes := consts.ResCreatedComment{
		Post:       fourthComment,
		FormErrors: fourthErrs,
		Error:      "",
	}

	fifthComment := newComment2
	fifthComment.Text = utils.RandomString(1400)
	fifthErrs := consts.EmptyCreateCommentErrors
	fifthErrs[0].Bool = true
	fifthErrs[0].Message = "This field must have less than 1275 characters"
	fifthReq, err := NewPostRequestComment(fifthComment, "/comments")
	fifthExpectedRes := consts.ResCreatedComment{
		Post:       fifthComment,
		FormErrors: fifthErrs,
		Error:      "",
	}

	sixthComment := newComment
	sixthComment.PostID = 0
	sixthComment.Text = ""
	sixthComment.ImageID = sql.NullInt32{
		Int32: 0,
		Valid: false,
	}
	sixthErrs := [2]consts.FormInputError{
		{
			Bool:    true,
			Message: "This field is required",
			Field:   "text",
		},
		{
			Bool:    false,
			Message: "",
			Field:   "image",
		},
	}
	pr6, pw6 := io.Pipe()
	form6 := multipart.NewWriter(pw6)
	go NewPostRequestCommentImage(t, pw6, form6, "test.png", sixthComment)

	sixthReq, err := http.NewRequest(http.MethodPost, "/posts", pr6)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	sixthReq.Header.Set("Content-Type", form6.FormDataContentType())
	sixthExpectedRes := consts.ResCreatedComment{
		Post:       sixthComment,
		FormErrors: sixthErrs,
		Error:      "",
	}

	pr7, pw7 := io.Pipe()
	form7 := multipart.NewWriter(pw7)
	go NewPostRequestCommentImage(t, pw7, form7, "textfile.txt", newComment)

	seventhReq, err := http.NewRequest(http.MethodPost, "/posts", pr7)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	seventhErrs := consts.EmptyCreateCommentErrors
	seventhPost := newComment
	seventhPost.PostID = 0
	seventhErrs[1].Bool = true
	seventhErrs[1].Message = "File type should be jpeg or png"
	seventhReq.Header.Set("Content-Type", form7.FormDataContentType())
	seventhExpectedRes := consts.ResCreatedComment{
		Post:       seventhPost,
		FormErrors: seventhErrs,
		Error:      "",
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
			Name:         "unsuccesful post without image, text too long",
			Req:          fifthReq,
			ExpectedRes:  fifthExpectedRes,
			ExpectedCode: http.StatusUnprocessableEntity,
			TestAfter: AfterRes{
				Valid: false,
				Type:  "",
			},
		},
		{
			Name:         "unsuccessful post post with image, text required",
			Req:          sixthReq,
			ExpectedRes:  sixthExpectedRes,
			ExpectedCode: http.StatusUnprocessableEntity,
			TestAfter: AfterRes{
				Valid: false,
				Type:  "",
			},
		},
		{
			Name:         "unsuccessful post post, file not png or jpeg",
			Req:          seventhReq,
			ExpectedRes:  seventhExpectedRes,
			ExpectedCode: http.StatusUnprocessableEntity,
			TestAfter: AfterRes{
				Valid: false,
				Type:  "",
			},
		},
	}

	TestPost(t, testCases, Th.CreateComment)
}

func TestDeletePost(t *testing.T) {
	newestPost, err := Th.q.GetNewestPost(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	firstReq := httptest.NewRequest(http.MethodDelete, "/posts/"+strconv.Itoa(int(newestPost.PostID)), nil)
	secondReq := httptest.NewRequest(http.MethodDelete, "/posts/"+strconv.Itoa(1000000), nil)
	thirdReq := httptest.NewRequest(http.MethodDelete, "/posts/"+strconv.Itoa(int(newestPost.PostID-1)), nil)

	testCases := []GetTestCase{
		{
			Name:         "successful delete of post with no image",
			Req:          firstReq,
			ExpectedRes:  "",
			ExpectedCode: http.StatusOK,
		},
		{
			Name: "failed request for non existing post",
			Req:  secondReq,
			ExpectedRes: consts.ErrorMessage{
				Msg: sql.ErrNoRows.Error(),
			},
			ExpectedCode: http.StatusNotFound,
		},
		{
			Name:         "successful delete of post with image",
			Req:          thirdReq,
			ExpectedRes:  "",
			ExpectedCode: http.StatusOK,
		},
	}

	TestGet(t, testCases, Th.DeletePost)
}
