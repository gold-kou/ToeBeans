package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httpContext "github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"

	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	testingHelper "github.com/gold-kou/ToeBeans/backend/testing"
	"github.com/gold-kou/ToeBeans/backend/testing/dummy"
	"github.com/stretchr/testify/assert"
)

var successRegisterUserReq = `
{
  "email": "testUser1@example.com",
  "password": "Password1234"
}
`

var errorRespWithoutUserName = `
{
  "status": 400,
  "message": "cannot be blank"
}
`

var errorRegisterUserReqWithoutEmail = `
{
  "password": "Password1234"
}
`
var errorRegisterUserRespWithoutEmail = `
{
  "status": 400,
  "message": "email: cannot be blank."
}
`

var errorRegisterUserReqWithoutPassword = `
{
  "email": "testUser1@example.com"
}
`
var errorRegisterUserRespWithoutPassword = `
{
  "status": 400,
  "message":"password: cannot be blank."
}
`

var errorRegisterUserReqNotEmailFormat = `
{
  "email": "hoge",
  "password": "Password1234"
}
`
var errorRegisterUserRespNotEmailFormat = `
{
  "status": 400,
  "message": "email: must be a valid email address."
}
`

var errorRegisterUserUserNameShort = `
{
  "email": "testUser1@example.com",
  "password": "Password1234"
}
`
var errorRespUserNameShort = `
{
  "status": 400,
  "message": "the length must be between 2 and 255"
}
`

var errReqRegisterUserNameNotAlphanumeric = `
{
  "email": "testUser1@example.com",
  "password": "Password1234"
}
`

var errRespUserNameNotAlphanumeric = `
{
  "status": 400,
  "message": "must contain English letters and digits only"
}
`

var errorRegisterUserPasswordShort = `
{
  "email": "testUser1@example.com",
  "password": "12345"
}
`
var errorRespPassword = `
{
  "status": 400,
  "message": "password: Your password must be at least 8 characters long, contain at least one number and have a mixture of uppercase and lowercase letters."
}
`

func TestRegisterUser(t *testing.T) {
	type args struct {
		userName string
		reqBody  string
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
			args:       args{userName: dummy.User1.Name, reqBody: successRegisterUserReq},
			method:     http.MethodPost,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error path parameter without user_name",
			args:       args{reqBody: successRegisterUserReq},
			method:     http.MethodPost,
			want:       errorRespWithoutUserName,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body without email",
			args:       args{userName: dummy.User1.Name, reqBody: errorRegisterUserReqWithoutEmail},
			method:     http.MethodPost,
			want:       errorRegisterUserRespWithoutEmail,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body without password",
			args:       args{userName: dummy.User1.Name, reqBody: errorRegisterUserReqWithoutPassword},
			method:     http.MethodPost,
			want:       errorRegisterUserRespWithoutPassword,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body not email format",
			args:       args{userName: dummy.User1.Name, reqBody: errorRegisterUserReqNotEmailFormat},
			method:     http.MethodPost,
			want:       errorRegisterUserRespNotEmailFormat,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error path parameter user_name short",
			args:       args{userName: "1", reqBody: errorRegisterUserUserNameShort},
			method:     http.MethodPost,
			want:       errorRespUserNameShort,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body user_name is not alphanumeric",
			args:       args{userName: "test_1", reqBody: errReqRegisterUserNameNotAlphanumeric},
			method:     http.MethodPost,
			want:       errRespUserNameNotAlphanumeric,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body password short",
			args:       args{userName: dummy.User1.Name, reqBody: errorRegisterUserPasswordShort},
			method:     http.MethodPost,
			want:       errorRespPassword,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not allowed method",
			args:       args{userName: dummy.User1.Name, reqBody: errorRegisterUserPasswordShort},
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

			// http request
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/users/%s", tt.args.userName), strings.NewReader(tt.args.reqBody))
			assert.NoError(t, err)
			vars := map[string]string{"user_name": tt.args.userName}
			req = mux.SetURLVars(req, vars)
			req.Header.Add(helper.HeaderKeyContentType, "application/json")
			resp := httptest.NewRecorder()

			// test target
			UserController(resp, req)
			assert.NoError(t, err)

			// assert db
			if tt.wantStatus == 200 {
				users, err := testingHelper.FindAllUsers(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, 1, len(users))
				// not be able to check password
				assert.Equal(t, users[0].Name, dummy.User1.Name)
				assert.Equal(t, users[0].Email, dummy.User1.Email)
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

var successRespGetUser = `
{
  "user_name": "testUser1",
  "icon": "UNKNOWN",
  "self_introduction": "UNKNOWN",
  "posting_count": 1,
  "like_count": 1,
  "liked_count": 1,
  "follow_count": 1,
  "followed_count": 1,
  "created_at": "2020-01-01T00:00:00+09:00"
}
`
var errorRespGetUserNameShort = `
{
  "status": 400,
  "message": "the length must be between 2 and 255"
}
`
var errRespGetUserNameNotAlphanumeric = `
{
  "status": 400,
  "message": "must contain English letters and digits only"
}
`
var errRespGetUserNotExists = `
{
  "status": 404,
  "message": "not exists data error"
}
`

func TestGetUser(t *testing.T) {
	type args struct {
		userName string
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
			args:       args{userName: dummy.User1.Name},
			method:     http.MethodGet,
			want:       successRespGetUser,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error user_name short",
			args:       args{userName: "a"},
			method:     http.MethodGet,
			want:       errorRespGetUserNameShort,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error user_name is not alphanumeric",
			args:       args{userName: "test_user"},
			method:     http.MethodGet,
			want:       errRespGetUserNameNotAlphanumeric,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error not exists user_name",
			args:       args{userName: "testUser0"},
			method:     http.MethodGet,
			want:       errRespGetUserNotExists,
			wantStatus: http.StatusNotFound,
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
			likeRepo := repository.NewLikeRepository(db)
			followRepo := repository.NewFollowRepository(db)
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)
			err = userRepo.Create(context.Background(), &dummy.User2)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting1)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting2)
			assert.NoError(t, err)
			err = likeRepo.Create(context.Background(), &dummy.Like1to2)
			assert.NoError(t, err)
			err = likeRepo.Create(context.Background(), &dummy.Like2to1)
			assert.NoError(t, err)
			err = followRepo.Create(context.Background(), &dummy.Follow1to2)
			assert.NoError(t, err)
			err = followRepo.Create(context.Background(), &dummy.Follow2to1)
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/users?user_name=%s", tt.args.userName), nil)
			assert.NoError(t, err)
			req = req.WithContext(httpContext.SetTokenUserName(req.Context(), dummy.User1.Name))
			resp := httptest.NewRecorder()

			// test target
			UserController(resp, req)
			assert.NoError(t, err)

			// assert http
			assert.Equal(t, tt.wantStatus, resp.Code)
			respBodyByte, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			if tt.wantStatus == http.StatusOK {
				var respGetUser modelHTTP.ResponseGetUser
				if err = json.Unmarshal(respBodyByte, &respGetUser); err != nil {
					assert.NoError(t, err)
				}
				respGetUser.CreatedAt = lib.NowFunc()
				respGetUserJSON, err := json.Marshal(respGetUser)
				assert.NoError(t, err)
				assert.JSONEq(t, tt.want, string(respGetUserJSON))
			} else {
				respBody := string(respBodyByte)
				assert.JSONEq(t, tt.want, respBody)
			}
		})
	}
}

var successReqUpdateUser = `
{
  "password": "NewPassword1234",
  "icon": "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",
  "self_introduction": "Hello!"
}
`
var errReqUpdateUserInvalidPassword = `
{
  "password": "password"
}
`
var errReqUpdateUserShortSelfIntroduction = `
{
  "self_introduction": "a"
}
`
var errReqUpdateUserIconDecodeFail = `
{
  "icon": "a"
}
`

var errRespUpdateUserInvalidPassword = `
{
  "status": 400,
  "message": "Your password must be at least 8 characters long, contain at least one number and have a mixture of uppercase and lowercase letters"
}
`
var errRespUpdateUserSelfIntroductionShort = `
{
  "status": 400,
  "message": "the length must be between 2 and 255"
}
`
var errRespUpdateUserDecodeFailure = `
{
  "status": 400,
  "message": "image decode failure"
}
`

func TestUpdateUser(t *testing.T) {
	type args struct {
		userName string
		reqBody  string
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
			args:       args{userName: dummy.User1.Name, reqBody: successReqUpdateUser},
			method:     http.MethodPut,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error password validation error",
			args:       args{userName: dummy.User1.Name, reqBody: errReqUpdateUserInvalidPassword},
			method:     http.MethodPut,
			want:       errRespUpdateUserInvalidPassword,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error self-introduction short",
			args:       args{userName: dummy.User1.Name, reqBody: errReqUpdateUserShortSelfIntroduction},
			method:     http.MethodPut,
			want:       errRespUpdateUserSelfIntroductionShort,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error icon decode failure",
			args:       args{userName: dummy.User1.Name, reqBody: errReqUpdateUserIconDecodeFail},
			method:     http.MethodPut,
			want:       errRespUpdateUserDecodeFailure,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error forbidden guest user",
			args:       args{userName: helper.GuestUserName, reqBody: successReqUpdateUser},
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
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/users/%s", tt.args.userName), strings.NewReader(tt.args.reqBody))
			assert.NoError(t, err)
			vars := map[string]string{"user_name": tt.args.userName}
			req = mux.SetURLVars(req, vars)
			if tt.name == "error forbidden guest user" {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), helper.GuestUserName))
			} else {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), dummy.User1.Name))
			}
			resp := httptest.NewRecorder()

			// test target
			UserController(resp, req)
			assert.NoError(t, err)

			if tt.wantStatus == 200 {
				users, err := testingHelper.FindAllUsers(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, "http://localhost:9000/toebeans-icons/testUser1", users[0].Icon)
				assert.Equal(t, "Hello!", users[0].SelfIntroduction)
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

var errRespDeleteUserNotExistsUserNameSpecified = `
{
  "message":"the user doesn't exist",
  "status":409
}
`

func TestDeleteUser(t *testing.T) {
	type args struct {
		userName string
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
			args:       args{userName: dummy.User1.Name},
			method:     http.MethodDelete,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error not existing user is specified",
			args:       args{userName: "notExistingUser"},
			method:     http.MethodDelete,
			want:       errRespDeleteUserNotExistsUserNameSpecified,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "error forbidden guest user",
			args:       args{userName: dummy.User1.Name},
			method:     http.MethodDelete,
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
			err = userRepo.Create(context.Background(), &dummy.User2)
			assert.NoError(t, err)
			postingRepo := repository.NewPostingRepository(db)
			err = postingRepo.Create(context.Background(), &dummy.Posting1)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting2)
			assert.NoError(t, err)
			followRepo := repository.NewFollowRepository(db)
			err = followRepo.Create(context.Background(), &dummy.Follow1to2)
			assert.NoError(t, err)
			err = followRepo.Create(context.Background(), &dummy.Follow2to1)
			assert.NoError(t, err)
			commentRepo := repository.NewCommentRepository(db)
			err = commentRepo.Create(context.Background(), &dummy.Comment1)
			assert.NoError(t, err)
			likeRepo := repository.NewLikeRepository(db)
			err = likeRepo.Create(context.Background(), &dummy.Like1to2)
			assert.NoError(t, err)
			err = likeRepo.Create(context.Background(), &dummy.Like2to1)
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/users/%s", tt.args.userName), nil)
			assert.NoError(t, err)
			vars := map[string]string{"user_name": tt.args.userName}
			req = mux.SetURLVars(req, vars)
			if tt.name == "error forbidden guest user" {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), helper.GuestUserName))
			} else {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), tt.args.userName))
			}
			resp := httptest.NewRecorder()

			// test target
			UserController(resp, req)
			assert.NoError(t, err)

			// assert db
			if tt.wantStatus == 200 {
				users, err := testingHelper.FindAllUsers(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, 1, len(users)) // user2がいるため
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

var successRespUserActivationHTML = `<!DOCTYPE html>
<html lang='en'>
<head>
    <meta charset='UTF-8'>
    <meta name='viewport' content='width=device-width, initial-scale=1.0'>
    <title>User Activation</title>
</head>
<body>
    <h1>Welcome to ToeBeans!</h1>
    User activation is success!
	<br>
	<br>
	<a href="http://localhost:3000/login">Login Page</a>
</body>
</html>`

var errRespUserActivationWithoutUserName = `
{
  "status": 400,
  "message": "cannot be blank"
}
`
var errRespUserActivationWithoutActivationKey = `
{
  "status": 400,
  "message": "cannot be blank"
}
`
var errRespUserActivationNameShort = `
{
  "status": 400,
  "message": "the length must be between 2 and 255"
}
`
var errRespUserActivationNameNotAlphanumeric = `
{
  "status": 400,
  "message": "must contain English letters and digits only"
}
`
var errRespUserActivationActivationKeyNotUUID = `
{
  "status": 400,
  "message": "must be a valid UUID"
}
`
var errRespUserActivationNotFound = `
{
  "status": 400,
  "message": "wrong user_name or activation_key, or might be already activated"
}
`

func TestUserActivation(t *testing.T) {
	type args struct {
		userName      string
		activationKey string
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
			args:       args{userName: dummy.User1.Name, activationKey: dummy.User1.ActivationKey},
			method:     http.MethodGet,
			want:       successRespUserActivationHTML,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error user_name is empty",
			args:       args{activationKey: dummy.User1.ActivationKey},
			method:     http.MethodGet,
			want:       errRespUserActivationWithoutUserName,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error user_name is short",
			args:       args{userName: "a", activationKey: dummy.User1.ActivationKey},
			method:     http.MethodGet,
			want:       errRespUserActivationNameShort,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error user_name is not alphanumeric",
			args:       args{userName: "test_user1", activationKey: dummy.User1.ActivationKey},
			method:     http.MethodGet,
			want:       errRespUserActivationNameNotAlphanumeric,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error activation_key is empty",
			args:       args{userName: dummy.User1.Name},
			method:     http.MethodGet,
			want:       errRespUserActivationWithoutActivationKey,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error activation_key not UUID",
			args:       args{userName: dummy.User1.Name, activationKey: "abc"},
			method:     http.MethodGet,
			want:       errRespUserActivationActivationKeyNotUUID,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error not existing user_name",
			args:       args{userName: "testUser0", activationKey: dummy.User1.ActivationKey},
			method:     http.MethodGet,
			want:       errRespUserActivationNotFound,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error not existing activation_key",
			args:       args{userName: dummy.User1.Name, activationKey: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"},
			method:     http.MethodGet,
			want:       errRespUserActivationNotFound,
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
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/user-activation/%s/%s", tt.args.userName, tt.args.activationKey), nil)
			assert.NoError(t, err)
			vars := map[string]string{"user_name": tt.args.userName, "activation_key": tt.args.activationKey}
			req = mux.SetURLVars(req, vars)
			resp := httptest.NewRecorder()

			// test target
			UserController(resp, req)
			assert.NoError(t, err)

			// assert db
			if tt.wantStatus == 200 {
				users, err := testingHelper.FindAllUsers(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, true, users[0].EmailVerified)
			}

			// assert http
			assert.Equal(t, tt.wantStatus, resp.Code)
			respBodyByte, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			respBody := string(respBodyByte)
			if tt.wantStatus == 200 {
				assert.Equal(t, tt.want, respBody)
			} else {
				assert.JSONEq(t, tt.want, respBody)
			}
		})
	}
}
