package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"time"

	"github.com/jobutterfly/olives/consts"
)

type PathInfo struct {
	Id int
}

const maxFileSize int64 = (1 << 20) / 2 // half a meg

func GetPathValues(ps []string, offset int) (PathInfo, error) {
	r := PathInfo{
		Id: 0,
	}

	if len(ps) > 3+offset {
		if ps[3+offset] != "" {
			err := consts.PathNotFound
			return r, err
		}
	}

	id, err := strconv.Atoi(ps[2+offset])
	if err != nil {
		err := consts.PathNotAnInteger
		return r, err
	}
	r.Id = id

	return r, err
}

func ValidateNewUser(email string, username string, password string) (errs [3]consts.FormInputError, valid bool) {
	valid = true
	errs = consts.EmptyCreateUserErrors

	_, err := mail.ParseAddress(email)
	if err != nil {
		errs[0].Bool = true
		errs[0].Message = "Invalid email address"
		valid = false
	} else if len(email) > 255 {
		errs[0].Bool = true
		errs[0].Message = "This field must have less than 275 characters"
		valid = false
	}

	if username == "" {
		errs[1].Bool = true
		errs[1].Message = "This field is required"
		valid = false
	} else if len(username) < 5 {
		errs[1].Bool = true
		errs[1].Message = "This field must be greater than 6 characters"
		valid = false
	} else if len(username) > 255 {
		errs[1].Bool = true
		errs[1].Message = "This field must have less than 275 characters"
		valid = false
	}

	if password == "" {
		errs[2].Bool = true
		errs[2].Message = "This field is required"
		valid = false
	} else if len(password) < 5 {
		errs[2].Bool = true
		errs[2].Message = "This field must be greater than 6 characters"
		valid = false
	} else if len(password) > 255 {
		errs[2].Bool = true
		errs[2].Message = "This field must have less than 275 characters"
		valid = false
	}

	return errs, valid
}

func ValidateNewComment(text string) (errs [2]consts.FormInputError, valid bool) {
	valid = true
	errs = consts.EmptyCreateCommentErrors

	if text == "" {
		errs[0].Bool = true
		errs[0].Message = "This field is required"
		valid = false
	} else if len(text) < 5 {
		errs[0].Bool = true
		errs[0].Message = "This field must be greater than 6 characters"
		valid = false
	} else if len(text) > 1275 {
		errs[0].Bool = true
		errs[0].Message = "This field must have less than 1275 characters"
		valid = false
	}

	return errs, valid
}

func ValidateNewCommentWithImage(text string, image multipart.File, header *multipart.FileHeader) (errs [2]consts.FormInputError, valid bool, imgPath string) {
	errs, valid = ValidateNewComment(text)
	imgPath = ""

	path, err := DownloadImage(image, header)
	if err != nil {
		errs[1].Bool = true
		errs[1].Message = err.Error()
		valid = false
		return errs, valid, imgPath
	}
	imgPath = path

	return errs, valid, imgPath
}

func ValidateNewPost(title string, text string) (errs [3]consts.FormInputError, valid bool) {
	valid = true
	errs = consts.EmptyCreatePostErrors

	if title == "" {
		errs[0].Bool = true
		errs[0].Message = "This field is required"
		valid = false
	} else if len(title) < 5 {
		errs[0].Bool = true
		errs[0].Message = "This field must be greater than 6 characters"
		valid = false
	} else if len(title) > 255 {
		errs[0].Bool = true
		errs[0].Message = "This field must have less than 255 characters"
		valid = false
	}

	if text == "" {
		errs[1].Bool = true
		errs[1].Message = "This field is required"
		valid = false
	} else if len(text) < 5 {
		errs[1].Bool = true
		errs[1].Message = "This field must be greater than 6 characters"
		valid = false
	} else if len(text) > 1275 {
		errs[1].Bool = true
		errs[1].Message = "This field must have less than 1275 characters"
		valid = false
	}

	return errs, valid
}

func ValidateNewPostWithImage(title string, text string, image multipart.File, header *multipart.FileHeader) (errs [3]consts.FormInputError, valid bool, imgPath string) {
	errs, valid = ValidateNewPost(title, text)
	imgPath = ""

	path, err := DownloadImage(image, header)
	if err != nil {
		errs[2].Bool = true
		errs[2].Message = err.Error()
		valid = false
		return errs, valid, imgPath
	}
	imgPath = path

	return errs, valid, imgPath
}

func DownloadImage(image multipart.File, header *multipart.FileHeader) (string, error) {
	if header.Size > maxFileSize {
		return "", errors.New("File size greater than 512 kilobytes. Choose a smaller file.")
	}

	tBuf := make([]byte, 512)
	if _, err := image.Read(tBuf); err != nil {
		return "", errors.New("Error when checking file type")
	}
	contentType := http.DetectContentType(tBuf)
	if contentType != "image/png" && contentType != "image/jpeg" {
		return "", errors.New("File type should be jpeg or png")
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, image); err != nil {
		return "", err
	}
	path := fmt.Sprintf("../view/images/%d%s", time.Now().Unix(), header.Filename)

	if err := os.WriteFile(path, buf.Bytes(), 0666); err != nil {
		return "", err
	}

	return path, nil
}
