package controller

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	testingHelper "github.com/gold-kou/ToeBeans/backend/testing"
	"github.com/gold-kou/ToeBeans/backend/testing/dummy"
	"github.com/stretchr/testify/assert"
)

var successReqPasswordChange = `
{
  "old_password": "Password1234",
  "new_password": "Password5678"
}
`
var errReqPasswordChangeNoOldPassword = `
{
  "new_password": "Password5678"
}
`
var errReqPasswordChangeNoNewPassword = `
{
  "old_password": "Password1234"
}
`
var errReqPasswordChangeInvalidPassword = `
{
  "old_password": "password",
  "new_password": "password"
}
`
var errReqPasswordChangeWrongOldPassword = `
{
  "old_password": "Hoge1234",
  "new_password": "Password5678"
}
`

var errRespPasswordChangeNoOldPassword = `
{
  "status": 400,
  "message": "old_password: cannot be blank."
}
`
var errRespPasswordChangeNoNewPassword = `
{
  "status": 400,
  "message": "new_password: cannot be blank."
}
`
var errRespPasswordChangeInvalidPassword = `
{
  "status": 400,
  "message": "new_password: Your password must be at least 8 characters long, contain at least one number and have a mixture of uppercase and lowercase letters; old_password: Your password must be at least 8 characters long, contain at least one number and have a mixture of uppercase and lowercase letters."
}
`
var errRespPasswordChangeWrongOldPassword = `
{
  "status": 400,
  "message": "not correct password"
}
`

func TestPasswordChange(t *testing.T) {
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
		//{
		//	name:       "success",
		//	args:       args{reqBody: successReqPasswordChange},
		//	method:     http.MethodPut,
		//	want:       testingHelper.RespSimpleSuccess,
		//	wantStatus: http.StatusOK,
		//},
		{
			name:       "error no old password",
			args:       args{reqBody: errReqPasswordChangeNoOldPassword},
			method:     http.MethodPut,
			want:       errRespPasswordChangeNoOldPassword,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error no new password",
			args:       args{reqBody: errReqPasswordChangeNoNewPassword},
			method:     http.MethodPut,
			want:       errRespPasswordChangeNoNewPassword,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error invalid password",
			args:       args{reqBody: errReqPasswordChangeInvalidPassword},
			method:     http.MethodPut,
			want:       errRespPasswordChangeInvalidPassword,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error old password is wrong",
			args:       args{reqBody: errReqPasswordChangeWrongOldPassword},
			method:     http.MethodPut,
			want:       errRespPasswordChangeWrongOldPassword,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error forbidden guest user",
			args:       args{reqBody: successReqPasswordChange},
			method:     http.MethodPut,
			want:       testingHelper.ErrForbidden,
			wantStatus: http.StatusForbidden,
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

			// http request
			req, err := http.NewRequest(tt.method, "/password", strings.NewReader(tt.args.reqBody))
			assert.NoError(t, err)
			var token string
			if tt.name == "error forbidden guest user" {
				token, err = lib.GenerateToken(lib.GuestUserName)
			} else {
				token, err = lib.GenerateToken(dummy.User1.Name)
			}
			assert.NoError(t, err)
			cookie := &http.Cookie{
				Name:  helper.CookieIDToken,
				Value: token,
			}
			req.AddCookie(cookie)
			resp := httptest.NewRecorder()

			// test target
			PasswordController(resp, req)
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
