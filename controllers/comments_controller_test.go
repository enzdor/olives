package controllers

import (
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

func TestCreateComment(t *testing.T) {
	if err := Start(); err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	newestComment, err := Th.q.GetNewestComment(context.Background())
	if err != nil {
		t.Errorf("expected no errors, got %v", err)
		return
	}

	newComment := utils.RandomComment()
	newComment.CommentID = newestComment.CommentID + 1

	firstExpectedRes := consts.ResCreatedComment{
		Comment:    newComment,
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
		PostID: 10,
	}
	secondExpectedRes := consts.ResCreatedComment{
		Comment: newComment2,
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
	firstReq, err := http.NewRequest(http.MethodPost, "/comments/create", pr)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	firstReq.Header.Set("Content-Type", form.FormDataContentType())

	pr2, pw2 := io.Pipe()
	form2 := multipart.NewWriter(pw2)
	go NewPostRequestCommentImage(t, pw2, form2, "cheese.png", newComment2)
	secondReq, err := http.NewRequest(http.MethodPost, "/comments/create", pr2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	secondReq.Header.Set("Content-Type", form2.FormDataContentType())

	thirdComment := utils.RandomComment()
	thirdComment.Text = "this is a valid text"
	thirdReq, err := NewPostRequestComment(thirdComment, "/comments/create")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	thirdExpectedRes := consts.ResCreatedComment{
		Comment:    thirdComment,
		FormErrors: consts.EmptyCreateCommentErrors,
		Error:      "",
	}

	fourthComment := newComment2
	fourthComment.Text = "bla"
	fourthErrs := consts.EmptyCreateCommentErrors
	fourthErrs[0].Bool = true
	fourthErrs[0].Message = "This field must be greater than 6 characters"
	fourthReq, err := NewPostRequestComment(fourthComment, "/comments/create")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	fourthExpectedRes := consts.ResCreatedComment{
		Comment:    fourthComment,
		FormErrors: fourthErrs,
		Error:      "",
	}

	fifthComment := newComment2
	fifthComment.Text = utils.RandomString(1300)
	fifthErrs := consts.EmptyCreateCommentErrors
	fifthErrs[0].Bool = true
	fifthErrs[0].Message = "This field must have less than 1275 characters"
	fifthReq, err := NewPostRequestComment(fifthComment, "/comments/create")
	fifthExpectedRes := consts.ResCreatedComment{
		Comment:    fifthComment,
		FormErrors: fifthErrs,
		Error:      "",
	}

	sixthComment := newComment
	sixthComment.CommentID = 0
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

	sixthReq, err := http.NewRequest(http.MethodPost, "/comments/create", pr6)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	sixthReq.Header.Set("Content-Type", form6.FormDataContentType())
	sixthExpectedRes := consts.ResCreatedComment{
		Comment:    sixthComment,
		FormErrors: sixthErrs,
		Error:      "",
	}

	pr7, pw7 := io.Pipe()
	form7 := multipart.NewWriter(pw7)
	go NewPostRequestCommentImage(t, pw7, form7, "textfile.txt", newComment)

	seventhReq, err := http.NewRequest(http.MethodPost, "/comment/creates", pr7)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	seventhErrs := consts.EmptyCreateCommentErrors
	seventhErrs[1].Bool = true
	seventhComment := newComment
	seventhComment.CommentID = 0
	seventhComment.ImageID.Int32 = 0
	seventhComment.ImageID.Valid = false
	seventhErrs[1].Message = "File type should be jpeg or png"
	seventhReq.Header.Set("Content-Type", form7.FormDataContentType())
	seventhExpectedRes := consts.ResCreatedComment{
		Comment:    seventhComment,
		FormErrors: seventhErrs,
		Error:      "",
	}

	testCases := []PostTestCase{
		{
			Name:         "successful post comment",
			Req:          firstReq,
			ExpectedRes:  firstExpectedRes,
			ExpectedCode: http.StatusCreated,
			TestAfter: AfterRes{
				Valid: true,
				Type:  "comment",
			},
		},
		{
			Name:         "unsuccessful post comment, image too large",
			Req:          secondReq,
			ExpectedRes:  secondExpectedRes,
			ExpectedCode: http.StatusUnprocessableEntity,
			TestAfter: AfterRes{
				Valid: false,
				Type:  "",
			},
		},
		{
			Name:         "succesful post comment without image",
			Req:          thirdReq,
			ExpectedRes:  thirdExpectedRes,
			ExpectedCode: http.StatusCreated,
			TestAfter: AfterRes{
				Valid: true,
				Type:  "comment",
			},
		},
		{
			Name:         "unsuccesful post comment  without image, text too short",
			Req:          fourthReq,
			ExpectedRes:  fourthExpectedRes,
			ExpectedCode: http.StatusUnprocessableEntity,
			TestAfter: AfterRes{
				Valid: false,
				Type:  "",
			},
		},
		{
			Name:         "unsuccesful post comment without image, text too long",
			Req:          fifthReq,
			ExpectedRes:  fifthExpectedRes,
			ExpectedCode: http.StatusUnprocessableEntity,
			TestAfter: AfterRes{
				Valid: false,
				Type:  "",
			},
		},
		{
			Name:         "unsuccessful post comment with image, text required",
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

func TestDeleteComment(t *testing.T) {
	if err := Start(); err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	newestComment, err := Th.q.GetNewestComment(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	firstReq := httptest.NewRequest(http.MethodDelete, "/comments/delete/"+strconv.Itoa(int(newestComment.CommentID)), nil)
	secondReq := httptest.NewRequest(http.MethodDelete, "/comments/delete/"+strconv.Itoa(1000000), nil)
	thirdReq := httptest.NewRequest(http.MethodDelete, "/comments/delete/"+strconv.Itoa(int(newestComment.CommentID-1)), nil)

	testCases := []GetTestCase{
		{
			Name:         "successful delete of comment with no image",
			Req:          firstReq,
			ExpectedRes:  "",
			ExpectedCode: http.StatusOK,
		},
		{
			Name: "failed request for non existing comment",
			Req:  secondReq,
			ExpectedRes: consts.ErrorMessage{
				Msg: sql.ErrNoRows.Error(),
			},
			ExpectedCode: http.StatusNotFound,
		},
		{
			Name:         "successful delete of comment with image",
			Req:          thirdReq,
			ExpectedRes:  "",
			ExpectedCode: http.StatusOK,
		},
	}

	TestGet(t, testCases, Th.DeleteComment)
}
