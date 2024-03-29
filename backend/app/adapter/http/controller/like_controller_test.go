package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	httpContext "github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"

	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"

	"github.com/gold-kou/ToeBeans/backend/app/lib"
	"github.com/gold-kou/ToeBeans/backend/testing/dummy"

	testingHelper "github.com/gold-kou/ToeBeans/backend/testing"
	"github.com/stretchr/testify/assert"
)

var errRespRegisterLikeWithoutPostingID = `
{
  "status": 400,
  "message": "cannot be blank"
}
`

var errRespRegisterLikeDuplicate = `
{
  "status": 409,
  "message": "Whoops, you already liked the posting"
}
`

var errRespRegisterLikeYourself = `
{
  "status": 409,
  "message": "you can't like your posting"
}
`

func TestRegisterLike(t *testing.T) {
	type args struct {
		postingID int64
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
			args:       args{postingID: dummy.Posting2.ID},
			method:     http.MethodPost,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error empty posting_id",
			args:       args{},
			method:     http.MethodPost,
			want:       errRespRegisterLikeWithoutPostingID,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:         "error duplicate like",
			args:         args{postingID: dummy.Posting2.ID},
			duplicateErr: true,
			method:       http.MethodPost,
			want:         errRespRegisterLikeDuplicate,
			wantStatus:   http.StatusConflict,
		},
		{
			name:         "error like yourself",
			args:         args{postingID: dummy.Posting1.ID},
			duplicateErr: true,
			method:       http.MethodPost,
			want:         errRespRegisterLikeYourself,
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
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/likes/%v", tt.args.postingID), nil)
			assert.NoError(t, err)
			vars := map[string]string{"posting_id": strconv.Itoa(int(tt.args.postingID))}
			req = mux.SetURLVars(req, vars)
			req = req.WithContext(httpContext.SetTokenUserID(req.Context(), dummy.User1.ID))
			req = req.WithContext(httpContext.SetTokenUserName(req.Context(), dummy.User1.Name))
			resp := httptest.NewRecorder()

			// test target
			LikeController(resp, req)
			assert.NoError(t, err)

			if tt.duplicateErr {
				// 2nd same request
				req, err := http.NewRequest(tt.method, fmt.Sprintf("/likes/%v", tt.args.postingID), nil)
				assert.NoError(t, err)
				vars := map[string]string{"posting_id": strconv.Itoa(int(tt.args.postingID))}
				req = mux.SetURLVars(req, vars)
				req = req.WithContext(httpContext.SetTokenUserID(req.Context(), dummy.User1.ID))
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), dummy.User1.Name))
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
			if tt.wantStatus == http.StatusOK {
				likes, err := testingHelper.FindAllLikes(context.Background(), db)
				assert.NoError(t, err)
				dummy.Like1to2.CreatedAt = lib.NowFunc()
				dummy.Like1to2.UpdatedAt = lib.NowFunc()
				likes[0].CreatedAt = lib.NowFunc()
				likes[0].UpdatedAt = lib.NowFunc()
				assert.Equal(t, dummy.Like1to2, likes[0])
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
  "message": "cannot be blank"
}
`
var errRespDeleteLikeNotExistingPostingID = `
{
  "status": 409,
  "message": "can't delete not existing like"
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
			postingRepo := repository.NewPostingRepository(db)
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)
			err = userRepo.Create(context.Background(), &dummy.User2)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting1)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting2)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/likes/%v", tt.args.postingID), nil)
			assert.NoError(t, err)
			vars := map[string]string{"posting_id": strconv.Itoa(int(tt.args.postingID))}
			req = mux.SetURLVars(req, vars)
			req = req.WithContext(httpContext.SetTokenUserID(req.Context(), dummy.User1.ID))
			req = req.WithContext(httpContext.SetTokenUserName(req.Context(), dummy.User1.Name))
			resp := httptest.NewRecorder()
			LikeController(resp, req)
			assert.NoError(t, err)

			// http request
			req, err = http.NewRequest(tt.method, fmt.Sprintf("/likes/%v", tt.args.postingID), nil)
			assert.NoError(t, err)
			vars = map[string]string{"posting_id": strconv.Itoa(int(tt.args.postingID))}
			req = mux.SetURLVars(req, vars)
			req = req.WithContext(httpContext.SetTokenUserID(req.Context(), dummy.User1.ID))
			req = req.WithContext(httpContext.SetTokenUserName(req.Context(), dummy.User1.Name))
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
			if tt.wantStatus == http.StatusOK {
				likes, err := testingHelper.FindAllLikes(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, 0, len(likes))
			}
		})
	}
}
