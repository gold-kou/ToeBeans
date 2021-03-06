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

	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/lib"
	"github.com/gold-kou/ToeBeans/backend/testing/dummy"

	testingHelper "github.com/gold-kou/ToeBeans/backend/testing"
	"github.com/stretchr/testify/assert"
)

var successReqRegisterLike = `
{
  "posting_id": 2
}
`
var errReqRegisterLikeWithoutPostingID = `
{
}
`
var successReqRegisterLikeYourself = `
{
  "posting_id": 1
}
`

var errRespRegisterLikeWithoutPostingID = `
{
  "status": 400,
  "message": "posting_id: cannot be blank."
}
`

var errRespRegisterLikeDuplicate = `
{
  "status": 400,
  "message": "Whoops, you already liked the posting"
}
`

var errRespRegisterLikeYourself = `
{
  "status": 400,
  "message": "you can't like your posting"
}
`

func TestRegisterLike(t *testing.T) {
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
			args:       args{reqBody: successReqRegisterLike},
			method:     http.MethodPost,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error empty posting_id",
			args:       args{reqBody: errReqRegisterLikeWithoutPostingID},
			method:     http.MethodPost,
			want:       errRespRegisterLikeWithoutPostingID,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:         "error duplicate like",
			args:         args{reqBody: successReqRegisterLike},
			duplicateErr: true,
			method:       http.MethodPost,
			want:         errRespRegisterLikeDuplicate,
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "error like yourself",
			args:         args{reqBody: successReqRegisterLikeYourself},
			duplicateErr: true,
			method:       http.MethodPost,
			want:         errRespRegisterLikeYourself,
			wantStatus:   http.StatusBadRequest,
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
			err = userRepo.Create(context.Background(), &dummy.User2)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting1)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting2)
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, "/like", strings.NewReader(tt.args.reqBody))
			assert.NoError(t, err)
			var token string
			token, err = lib.GenerateToken(dummy.User1.Name)
			assert.NoError(t, err)
			cookie := &http.Cookie{
				Name:  helper.CookieIDToken,
				Value: token,
			}
			req.AddCookie(cookie)
			resp := httptest.NewRecorder()

			// test target
			LikeController(resp, req)
			assert.NoError(t, err)

			if tt.duplicateErr {
				// 2nd same request
				req, err := http.NewRequest(tt.method, "/like", strings.NewReader(tt.args.reqBody))
				assert.NoError(t, err)
				token, err := lib.GenerateToken(dummy.User1.Name)
				assert.NoError(t, err)
				cookie := &http.Cookie{
					Name:  helper.CookieIDToken,
					Value: token,
				}
				req.AddCookie(cookie)
				resp := httptest.NewRecorder()
				LikeController(resp, req)
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
				likes, err := testingHelper.FindAllLikes(context.Background(), db)
				assert.NoError(t, err)
				dummy.Like1.CreatedAt = lib.NowFunc()
				dummy.Like1.UpdatedAt = lib.NowFunc()
				likes[0].CreatedAt = lib.NowFunc()
				likes[0].UpdatedAt = lib.NowFunc()
				assert.Equal(t, dummy.Like1, likes[0])

				// increment check
				users, err := testingHelper.FindAllUsers(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), users[0].LikeCount)
				assert.Equal(t, int64(1), users[1].LikedCount)
				postings, err := testingHelper.FindAllPostings(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), postings[1].LikedCount)
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

var errRespDeleteLikeWithoutPostingID = `
{
  "status": 400,
  "message": "posting_id: cannot be blank"
}
`
var errRespDeleteLikeNotExistingPostingID = `
{
  "status": 400,
  "message": "not exists data error"
}
`

func TestDeleteLike(t *testing.T) {
	type args struct {
		postingID int64
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
			args:       args{postingID: dummy.Posting2.ID},
			method:     http.MethodDelete,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error empty posting_id",
			args:       args{},
			method:     http.MethodDelete,
			want:       errRespDeleteLikeWithoutPostingID,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error not existing posting_id",
			args:       args{postingID: 99999},
			method:     http.MethodDelete,
			want:       errRespDeleteLikeNotExistingPostingID,
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
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)
			err = userRepo.Create(context.Background(), &dummy.User2)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting1)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting2)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/like", strings.NewReader(successReqRegisterLike))
			assert.NoError(t, err)
			token, err := lib.GenerateToken(dummy.User1.Name)
			assert.NoError(t, err)
			cookie := &http.Cookie{
				Name:  helper.CookieIDToken,
				Value: token,
			}
			req.AddCookie(cookie)
			resp := httptest.NewRecorder()
			LikeController(resp, req)
			assert.NoError(t, err)

			// http request
			req, err = http.NewRequest(tt.method, fmt.Sprintf("/like/%v", tt.args.postingID), nil)
			assert.NoError(t, err)
			vars := map[string]string{"posting_id": strconv.Itoa(int(tt.args.postingID))}
			req = mux.SetURLVars(req, vars)
			token, err = lib.GenerateToken(dummy.User1.Name)
			assert.NoError(t, err)
			cookie = &http.Cookie{
				Name:  helper.CookieIDToken,
				Value: token,
			}
			req.AddCookie(cookie)
			resp = httptest.NewRecorder()

			// test target
			LikeController(resp, req)
			assert.NoError(t, err)

			// assert http
			assert.Equal(t, tt.wantStatus, resp.Code)
			respBodyByte, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			respBody := string(respBodyByte)
			assert.JSONEq(t, tt.want, respBody)

			// assert db
			if tt.wantStatus == 200 {
				likes, err := testingHelper.FindAllLikes(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, 0, len(likes))

				users, err := testingHelper.FindAllUsers(context.Background(), db)
				assert.NoError(t, err)
				for _, user := range users {
					if user.Name == dummy.User1.Name {
						// like
						assert.Equal(t, int64(0), users[0].LikeCount)
					}
					if user.Name == dummy.User2.Name {
						// liked
						assert.Equal(t, int64(0), users[1].LikedCount)
					}
				}

				postings, err := testingHelper.FindAllPostings(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, int64(0), postings[0].LikedCount)
			}
		})
	}
}
