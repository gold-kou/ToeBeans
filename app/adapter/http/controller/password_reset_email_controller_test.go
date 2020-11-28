package controller

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gold-kou/ToeBeans/app/domain/repository"
	"github.com/gold-kou/ToeBeans/testing/dummy"

	"github.com/gold-kou/ToeBeans/app/adapter/http/helper"
	testingHelper "github.com/gold-kou/ToeBeans/testing"
	"github.com/stretchr/testify/assert"
)

var successReqPasswordResetEmail = `
{
  "email": "testUser1@example.com"
}
`
var errReqPasswordResetEmailWithoutEmail = `
{
}
`
var errReqPasswordResetEmailNotEmailFormat = `
{
  "email": "abc"
}
`
var errReqPasswordResetEmailNotExistingEmail = `
{
  "email": "testUser0@example.com"
}
`

var errRespPasswordResetEmailWithoutEmail = `
{
  "status": 400,
  "message":"email: cannot be blank."
}
`
var errRespPasswordResetEmailNotEmailFormat = `
{
  "status": 400,
  "message":"email: must be a valid email address."
}
`
var errRespPasswordResetEmailNotExistingEmail = `
{
  "status": 400,
  "message":"not exists data error"
}
`
var errRespPasswordResetEmailOverLimitCount = `
{
  "status": 400,
  "message": "you can't reset password as it exceeds limit counts"
}
`

func TestPasswordResetEmailController(t *testing.T) {
	type args struct {
		reqBody string
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
			args:       args{reqBody: successReqPasswordResetEmail},
			method:     http.MethodPost,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error request body without email",
			args:       args{reqBody: errReqPasswordResetEmailWithoutEmail},
			method:     http.MethodPost,
			want:       errRespPasswordResetEmailWithoutEmail,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body not email format",
			args:       args{reqBody: errReqPasswordResetEmailNotEmailFormat},
			method:     http.MethodPost,
			want:       errRespPasswordResetEmailNotEmailFormat,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body not existing email",
			args:       args{reqBody: errReqPasswordResetEmailNotExistingEmail},
			method:     http.MethodPost,
			want:       errRespPasswordResetEmailNotExistingEmail,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error over password reset limit count",
			args:       args{reqBody: successReqPasswordResetEmail},
			method:     http.MethodPost,
			want:       errRespPasswordResetEmailOverLimitCount,
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
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)
			if tt.name == "error over password reset limit count" {
				err = testingHelper.UpdatePasswordResetEmailCount(db)
				assert.NoError(t, err)
			}

			// http request
			req, err := http.NewRequest(tt.method, "/password-reset-email", strings.NewReader(tt.args.reqBody))
			assert.NoError(t, err)
			req.Header.Add(helper.HeaderKeyContentType, "application/json")
			resp := httptest.NewRecorder()

			// test target
			PasswordResetEmailController(resp, req)
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
