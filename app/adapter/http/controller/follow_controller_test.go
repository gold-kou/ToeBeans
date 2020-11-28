package controller

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gold-kou/ToeBeans/app/domain/repository"

	"github.com/gold-kou/ToeBeans/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/app/lib"
	"github.com/gold-kou/ToeBeans/testing/dummy"

	testingHelper "github.com/gold-kou/ToeBeans/testing"
	"github.com/stretchr/testify/assert"
)

var successReqRegisterFollow = `
{
  "followed_user_name": "testUser2"
}
`

var errReqRegisterFollowWithoutUserName = `
{
}
`

var errReqRegisterFollowUserNameNotAlphanumeric = `
{
  "followed_user_name": "test_2"
}
`

var errRespRegisterFollowWithoutUserName = `
{
  "status": 400,
  "message": "followed_user_name: cannot be blank."
}
`

var errRespRegisterFollowUserNameNotAlphanumeric = `
{
  "status": 400,
  "message": "followed_user_name: must contain English letters and digits only."
}
`

var errRespRegisterFollowDuplicate = `
{
  "status": 400,
  "message": "Whoops, you already followed the user"
}
`

func TestRegisterFollow(t *testing.T) {
	type args struct {
		reqBody string
	}
	tests := []struct {
		name         string
		args         args
		duplicateErr bool
		method       string
		want         string
		wantStatus   int
	}{
		{
			name:       "success",
			args:       args{reqBody: successReqRegisterFollow},
			method:     http.MethodPost,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error empty user_name",
			args:       args{reqBody: errReqRegisterFollowWithoutUserName},
			method:     http.MethodPost,
			want:       errRespRegisterFollowWithoutUserName,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error user_name not alphanumeric",
			args:       args{reqBody: errReqRegisterFollowUserNameNotAlphanumeric},
			method:     http.MethodPost,
			want:       errRespRegisterFollowUserNameNotAlphanumeric,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:         "error duplicate follow",
			args:         args{reqBody: successReqRegisterFollow},
			duplicateErr: true,
			method:       http.MethodPost,
			want:         errRespRegisterFollowDuplicate,
			wantStatus:   http.StatusBadRequest,
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
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)
			err = userRepo.Create(context.Background(), &dummy.User2)
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, "/follow", strings.NewReader(tt.args.reqBody))
			assert.NoError(t, err)
			token, err := lib.GenerateToken(dummy.User1.Name)
			assert.NoError(t, err)
			req.Header.Add(helper.HeaderKeyAuthorization, "Bearer "+token)
			resp := httptest.NewRecorder()

			// test target
			FollowController(resp, req)
			assert.NoError(t, err)

			if tt.duplicateErr {
				// 2nd same request
				req, err := http.NewRequest(tt.method, "/follow", strings.NewReader(tt.args.reqBody))
				assert.NoError(t, err)
				token, err := lib.GenerateToken(dummy.User1.Name)
				assert.NoError(t, err)
				req.Header.Add(helper.HeaderKeyAuthorization, "Bearer "+token)
				resp := httptest.NewRecorder()
				FollowController(resp, req)
				assert.NoError(t, err)

				// assert http
				assert.Equal(t, tt.wantStatus, resp.Code)
				respBodyByte, err := ioutil.ReadAll(resp.Body)
				assert.NoError(t, err)
				respBody := string(respBodyByte)
				assert.JSONEq(t, tt.want, respBody)
				return
			}

			// assert db
			if tt.wantStatus == 200 {
				follows, err := testingHelper.FindAllFollows(context.Background(), db)
				assert.NoError(t, err)
				dummy.Follow1.CreatedAt = lib.NowFunc()
				dummy.Follow1.UpdatedAt = lib.NowFunc()
				follows[0].CreatedAt = lib.NowFunc()
				follows[0].UpdatedAt = lib.NowFunc()
				assert.Equal(t, dummy.Follow1, follows[0])
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
