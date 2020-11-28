package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/app/domain/repository"

	"github.com/gold-kou/ToeBeans/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/app/lib"
	"github.com/gold-kou/ToeBeans/testing/dummy"

	testingHelper "github.com/gold-kou/ToeBeans/testing"
	"github.com/stretchr/testify/assert"
)

var errRespDeleteCommentWithoutCommentID = `
{
  "status": 400,
  "message": "parameter comment_id is required"
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
			name:       "not allowed method",
			args:       args{commentID: dummy.Comment1.ID},
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
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting1)
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/comment/%v", tt.args.commentID), nil)
			assert.NoError(t, err)
			vars := map[string]string{"comment_id": strconv.Itoa(int(tt.args.commentID))}
			req = mux.SetURLVars(req, vars)
			token, err := lib.GenerateToken(dummy.User1.Name)
			assert.NoError(t, err)
			req.Header.Add(helper.HeaderKeyAuthorization, "Bearer "+token)
			resp := httptest.NewRecorder()

			// test target
			CommentCommentIDController(resp, req)
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
