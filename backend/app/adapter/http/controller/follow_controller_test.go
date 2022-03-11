package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	httpContext "github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"

	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"

	"github.com/gold-kou/ToeBeans/backend/testing/dummy"

	testingHelper "github.com/gold-kou/ToeBeans/backend/testing"
	"github.com/stretchr/testify/assert"
)

var errRespFollowWithoutUserName = `
{
  "status": 400,
  "message": "cannot be blank"
}
`

var errRespFollowUserNameNotAlphanumeric = `
{
  "status": 400,
  "message": "must contain English letters and digits only"
}
`

var errRespRegisterFollowDuplicate = `
{
  "status": 409,
  "message": "Whoops, you already followed the user"
}
`

func TestRegisterFollow(t *testing.T) {
	type args struct {
		followedUserName string
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
			args:       args{followedUserName: dummy.User2.Name},
			method:     http.MethodPost,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error empty user_name",
			args:       args{},
			method:     http.MethodPost,
			want:       errRespFollowWithoutUserName,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error user_name not alphanumeric",
			args:       args{followedUserName: "test_2"},
			method:     http.MethodPost,
			want:       errRespFollowUserNameNotAlphanumeric,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error forbidden guest user",
			args:       args{followedUserName: dummy.User2.Name},
			method:     http.MethodPost,
			want:       testingHelper.ErrForbidden,
			wantStatus: http.StatusForbidden,
		},
		{
			name:         "error duplicate follow",
			args:         args{followedUserName: dummy.User2.Name},
			duplicateErr: true,
			method:       http.MethodPost,
			want:         errRespRegisterFollowDuplicate,
			wantStatus:   http.StatusConflict,
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
			err = userRepo.Create(context.Background(), &dummy.User2)
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/follows/%s", tt.args.followedUserName), nil)
			assert.NoError(t, err)
			vars := map[string]string{"followed_user_name": tt.args.followedUserName}
			req = mux.SetURLVars(req, vars)
			if tt.name == "error forbidden guest user" {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), helper.GuestUserName))
			} else {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), dummy.User1.Name))
			}
			resp := httptest.NewRecorder()

			// test target
			FollowController(resp, req)
			assert.NoError(t, err)

			if tt.duplicateErr {
				// 2nd same request
				req, err := http.NewRequest(tt.method, fmt.Sprintf("/follows/%s", tt.args.followedUserName), nil)
				assert.NoError(t, err)
				vars := map[string]string{"followed_user_name": tt.args.followedUserName}
				req = mux.SetURLVars(req, vars)
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), dummy.User1.Name))
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
			//if tt.wantStatus == 200 {
			//	follows, err := testingHelper.FindAllFollows(context.Background(), db)
			//	assert.NoError(t, err)
			//	dummy.Follow1to2.CreatedAt = lib.NowFunc()
			//	dummy.Follow1to2.UpdatedAt = lib.NowFunc()
			//	follows[0].CreatedAt = lib.NowFunc()
			//	follows[0].UpdatedAt = lib.NowFunc()
			//	assert.Equal(t, dummy.Follow1to2, follows[0])
			//}

			// assert http
			assert.Equal(t, tt.wantStatus, resp.Code)
			respBodyByte, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			respBody := string(respBodyByte)
			assert.JSONEq(t, tt.want, respBody)
		})
	}
}

var successRespGetFollowStateTrue = `
{
  "is_follow": true
}
`

var successRespGetFollowStateFalse = `
{
  "is_follow": false
}
`

func TestGetFollowState(t *testing.T) {
	type args struct {
		followedUserName string
	}
	tests := []struct {
		name       string
		args       args
		method     string
		want       string
		wantStatus int
	}{
		{
			name:       "success true",
			args:       args{followedUserName: dummy.User2.Name},
			method:     http.MethodGet,
			want:       successRespGetFollowStateTrue,
			wantStatus: http.StatusOK,
		},
		{
			name:       "success false",
			args:       args{followedUserName: dummy.User2.Name},
			method:     http.MethodGet,
			want:       successRespGetFollowStateFalse,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error empty user_name",
			args:       args{},
			method:     http.MethodGet,
			want:       errRespFollowWithoutUserName,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error user_name not alphanumeric",
			args:       args{followedUserName: "test_2"},
			method:     http.MethodGet,
			want:       errRespFollowUserNameNotAlphanumeric,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error forbidden guest user",
			args:       args{followedUserName: dummy.User2.Name},
			method:     http.MethodGet,
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
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)
			err = userRepo.Create(context.Background(), &dummy.User2)
			assert.NoError(t, err)
			if tt.want == successRespGetFollowStateTrue {
				followRepo := repository.NewFollowRepository(db)
				err := followRepo.Create(context.Background(), &dummy.Follow1to2)
				assert.NoError(t, err)
			}

			// http request
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/follows/%s", tt.args.followedUserName), nil)
			assert.NoError(t, err)
			vars := map[string]string{"followed_user_name": tt.args.followedUserName}
			req = mux.SetURLVars(req, vars)
			if tt.name == "error forbidden guest user" {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), helper.GuestUserName))
			} else {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), dummy.User1.Name))
			}
			resp := httptest.NewRecorder()

			// test target
			FollowController(resp, req)
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

var errRespDeleteFollowWithoutUserName = `
{
  "status": 400,
  "message": "cannot be blank"
}
`
var errRespDeleteNotExistingFollowedUserName = `
{
  "status": 409,
  "message": "the user doesn't exist"
}
`
var errRespDeleteNotExistingFollow = `
{
  "status": 409,
  "message": "can't delete not existing follow"
}
`

func TestDeleteFollow(t *testing.T) {
	type args struct {
		followedUserName string
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
			args:       args{followedUserName: dummy.User2.Name},
			method:     http.MethodDelete,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error empty user_name",
			args:       args{},
			method:     http.MethodDelete,
			want:       errRespDeleteFollowWithoutUserName,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error forbidden guest user",
			args:       args{followedUserName: dummy.User2.Name},
			method:     http.MethodDelete,
			want:       testingHelper.ErrForbidden,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "error not existing followed user name",
			args:       args{followedUserName: "notExisitngUser"},
			method:     http.MethodDelete,
			want:       errRespDeleteNotExistingFollowedUserName,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "error not existing follow",
			args:       args{followedUserName: dummy.User3.Name},
			method:     http.MethodDelete,
			want:       errRespDeleteNotExistingFollow,
			wantStatus: http.StatusConflict,
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
			err = userRepo.Create(context.Background(), &dummy.User2)
			assert.NoError(t, err)
			err = userRepo.Create(context.Background(), &dummy.User3)
			assert.NoError(t, err)
			followRepo := repository.NewFollowRepository(db)
			err = followRepo.Create(context.Background(), &dummy.Follow1to2)
			assert.NoError(t, err)
			err = userRepo.UpdateFollowCount(context.Background(), dummy.User1.Name, true)
			assert.NoError(t, err)
			err = userRepo.UpdateFollowedCount(context.Background(), dummy.User2.Name, true)
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/follows/%s", tt.args.followedUserName), nil)
			assert.NoError(t, err)
			vars := map[string]string{"followed_user_name": tt.args.followedUserName}
			req = mux.SetURLVars(req, vars)
			if tt.name == "error forbidden guest user" {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), helper.GuestUserName))
			} else {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), dummy.User1.Name))
			}
			resp := httptest.NewRecorder()

			// test target
			FollowController(resp, req)
			assert.NoError(t, err)

			// assert db
			if tt.wantStatus == 200 {
				follows, err := testingHelper.FindAllFollows(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, 0, len(follows))

				users, err := testingHelper.FindAllUsers(context.Background(), db)
				assert.NoError(t, err)
				for _, user := range users {
					if user.Name == dummy.User1.Name {
						// following
						assert.Equal(t, int64(0), users[0].FollowCount)
					}
					if user.Name == dummy.User2.Name {
						// followed
						assert.Equal(t, int64(0), users[1].FollowedCount)
					}
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
