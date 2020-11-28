package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gold-kou/ToeBeans/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
	"github.com/gold-kou/ToeBeans/app/lib"
	testingHelper "github.com/gold-kou/ToeBeans/testing"
	"github.com/gold-kou/ToeBeans/testing/dummy"
	"github.com/stretchr/testify/assert"
)

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
			want:       errorNotAllowedMethod,
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
			token, err := lib.GenerateToken(dummy.User1.Name)
			assert.NoError(t, err)
			req.Header.Add(helper.HeaderKeyAuthorization, "Bearer "+token)
			resp := httptest.NewRecorder()

			// test target
			CommentsController(resp, req)
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
