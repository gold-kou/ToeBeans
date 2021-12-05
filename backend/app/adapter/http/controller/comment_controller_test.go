package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	httpContext "github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
	"github.com/gold-kou/ToeBeans/backend/app/lib"
	testingHelper "github.com/gold-kou/ToeBeans/backend/testing"
	"github.com/gold-kou/ToeBeans/backend/testing/dummy"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var successReqRegisterComment = `
{
  "comment": "test comment"
}
`

var errReqRegisterCommentWithoutComment = `
{
}
`

var errRespRegisterCommentWithoutPostingID = `
{
  "status": 400,
  "message": "cannot be blank"
}
`

var errRespRegisterCommentWithoutComment = `
{
  "status": 400,
  "message": "comment: cannot be blank."
}
`

func TestRegisterComment(t *testing.T) {
	type args struct {
		postingID int64
		reqBody   string
	}
	tests := []struct {
		name       string
		args       args
		method     string
		want       string
		wantStatus int
	}{
		{
			name:       "success",
			args:       args{postingID: dummy.Posting1.ID, reqBody: successReqRegisterComment},
			method:     http.MethodPost,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error empty posting_id",
			args:       args{reqBody: successReqRegisterComment},
			method:     http.MethodPost,
			want:       errRespRegisterCommentWithoutPostingID,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error empty comment",
			args:       args{postingID: dummy.Posting1.ID, reqBody: errReqRegisterCommentWithoutComment},
			method:     http.MethodPost,
			want:       errRespRegisterCommentWithoutComment,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error forbidden guest user",
			args:       args{postingID: dummy.Posting1.ID, reqBody: successReqRegisterComment},
			method:     http.MethodPost,
			want:       testingHelper.ErrForbidden,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "not allowed method",
			args:       args{},
			method:     http.MethodHead,
			want:       testingHelper.ErrNotAllowedMethod,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// init
			db := testingHelper.SetupDBTest()
			defer testingHelper.TeardownDBTest(db)
			testingHelper.SetTestTime()
			defer testingHelper.ResetTime()

			// insert dummy data
			userRepo := repository.NewUserRepository(db)
			postingRepo := repository.NewPostingRepository(db)
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting1)
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/comments/%v", tt.args.postingID), strings.NewReader(tt.args.reqBody))
			assert.NoError(t, err)
			vars := map[string]string{"posting_id": strconv.Itoa(int(tt.args.postingID))}
			req = mux.SetURLVars(req, vars)
			if tt.name == "error forbidden guest user" {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), lib.GuestUserName))
			} else {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), dummy.User1.Name))
			}
			resp := httptest.NewRecorder()

			// test target
			CommentController(resp, req)

			// assert db
			if tt.wantStatus == 200 {
				comments, err := testingHelper.FindAllComments(context.Background(), db)
				assert.NoError(t, err)
				dummy.Comment1.CreatedAt = lib.NowFunc()
				dummy.Comment1.UpdatedAt = lib.NowFunc()
				comments[0].CreatedAt = lib.NowFunc()
				comments[0].UpdatedAt = lib.NowFunc()
				assert.Equal(t, dummy.Comment1, comments[0])
			}

			// assert http
			assert.Equal(t, tt.wantStatus, resp.Code)
			respBodyByte, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			respBody := string(respBodyByte)
			assert.JSONEq(t, tt.want, respBody)
		})
	}
}

var successRespGetComments = `
{
  "posting_id": 1,
  "comments": [
    {
      "comment_id": 1,
      "user_name": "testUser1",
      "commented_at": "2020-01-01T00:00:00+09:00",
      "comment": "test comment"
    }
  ]
}
`
var successRespGetCommentsEmpty = `
{
}
`
var errRespGetCommentsWithoutPostingID = `
{
  "status": 400,
  "message": "posting_id: cannot be blank."
}
`

func TestGetComments(t *testing.T) {
	type args struct {
		postingID string
	}
	tests := []struct {
		name       string
		args       args
		method     string
		want       string
		wantStatus int
	}{
		{
			name:       "success",
			args:       args{postingID: "1"},
			method:     http.MethodGet,
			want:       successRespGetComments,
			wantStatus: http.StatusOK,
		},
		{
			name:       "success no comments",
			args:       args{postingID: "2"},
			method:     http.MethodGet,
			want:       successRespGetCommentsEmpty,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error empty posting_id",
			args:       args{},
			method:     http.MethodGet,
			want:       errRespGetCommentsWithoutPostingID,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not allowed method",
			args:       args{},
			method:     http.MethodHead,
			want:       testingHelper.ErrNotAllowedMethod,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// init
			db := testingHelper.SetupDBTest()
			defer testingHelper.TeardownDBTest(db)
			testingHelper.SetTestTime()
			defer testingHelper.ResetTime()

			// insert dummy data
			userRepo := repository.NewUserRepository(db)
			postingRepo := repository.NewPostingRepository(db)
			commentRepo := repository.NewCommentRepository(db)
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)
			err = userRepo.Create(context.Background(), &dummy.User2)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting1)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting2)
			assert.NoError(t, err)
			err = commentRepo.Create(context.Background(), &dummy.Comment1)
			assert.NoError(t, err)
			err = testingHelper.UpdateNow(db, "comments")
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/comments?posting_id=%v", tt.args.postingID), nil)
			assert.NoError(t, err)
			resp := httptest.NewRecorder()
			req = req.WithContext(httpContext.SetTokenUserName(req.Context(), dummy.User1.Name))

			// test target
			CommentController(resp, req)
			assert.NoError(t, err)

			// assert http
			assert.Equal(t, tt.wantStatus, resp.Code)
			respBodyByte, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			respBody := string(respBodyByte)
			assert.JSONEq(t, tt.want, respBody)
		})
	}
}

var errRespDeleteCommentWithoutCommentID = `
{
  "status": 400,
  "message": "cannot be blank"
}
`

func TestDeleteComment(t *testing.T) {
	type args struct {
		commentID int64
	}
	tests := []struct {
		name       string
		args       args
		method     string
		want       string
		wantStatus int
	}{
		{
			name:       "success",
			args:       args{commentID: dummy.Comment1.ID},
			method:     http.MethodDelete,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error empty comment_id",
			args:       args{},
			method:     http.MethodDelete,
			want:       errRespDeleteCommentWithoutCommentID,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error forbidden guest user",
			args:       args{commentID: dummy.Comment1.ID},
			method:     http.MethodDelete,
			want:       testingHelper.ErrForbidden,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "not allowed method",
			args:       args{commentID: dummy.Comment1.ID},
			method:     http.MethodHead,
			want:       testingHelper.ErrNotAllowedMethod,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// init
			db := testingHelper.SetupDBTest()
			defer testingHelper.TeardownDBTest(db)
			testingHelper.SetTestTime()
			defer testingHelper.ResetTime()

			// insert dummy data
			userRepo := repository.NewUserRepository(db)
			postingRepo := repository.NewPostingRepository(db)
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting1)
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/comments/%v", tt.args.commentID), nil)
			assert.NoError(t, err)
			vars := map[string]string{"comment_id": strconv.Itoa(int(tt.args.commentID))}
			req = mux.SetURLVars(req, vars)
			if tt.name == "error forbidden guest user" {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), lib.GuestUserName))
			} else {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), dummy.User1.Name))
			}
			resp := httptest.NewRecorder()

			// test target
			CommentController(resp, req)
			assert.NoError(t, err)

			// db check
			if tt.wantStatus == 200 {
				comments, err := testingHelper.FindAllComments(context.Background(), db)
				assert.NoError(t, err)
				if len(comments) != 0 {
					t.Errorf("want is empty, but got %+v", comments)
				}
			}

			// assert http
			assert.Equal(t, tt.wantStatus, resp.Code)
			respBodyByte, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			respBody := string(respBodyByte)
			assert.JSONEq(t, tt.want, respBody)
		})
	}
}
