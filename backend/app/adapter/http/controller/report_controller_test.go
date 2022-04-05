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

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
	testingHelper "github.com/gold-kou/ToeBeans/backend/testing"
	"github.com/gold-kou/ToeBeans/backend/testing/dummy"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var successSubmitUserReportReq = `
{
  "detail": "an inappropriate user"
}
`
var errorSubmitUserReportRespWithoutUserID = `
{
  "status": 400,
  "message": "cannot be blank"
}
`
var errorSubmitUserReportReqWithoutDetail = `
{
}
`
var errorSubmitUserReportRespWithoutDetail = `
{
  "status": 400,
  "message": "detail: cannot be blank."
}
`
var errorSubmitUserReportRespNotExistsUser = `
{
  "status": 400,
  "message": "the user doesn't exist"
}
`

func TestSubmitUserReport(t *testing.T) {
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
			args:       args{userName: dummy.UserReport1.UserName, reqBody: successSubmitUserReportReq},
			method:     http.MethodPost,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error path parameter without user_id",
			args:       args{reqBody: successSubmitUserReportReq},
			method:     http.MethodPost,
			want:       errorSubmitUserReportRespWithoutUserID,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body without detail",
			args:       args{userName: dummy.UserReport1.UserName, reqBody: errorSubmitUserReportReqWithoutDetail},
			method:     http.MethodPost,
			want:       errorSubmitUserReportRespWithoutDetail,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error not exists user error",
			args:       args{userName: "not exists user name", reqBody: successSubmitPostingReportReq},
			method:     http.MethodPost,
			want:       errorSubmitUserReportRespNotExistsUser,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not allowed method",
			args:       args{userName: dummy.UserReport1.UserName, reqBody: successSubmitUserReportReq},
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
			userRepo := repository.NewUserRepository(db)
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/reports/users/%s", tt.args.userName), strings.NewReader(tt.args.reqBody))
			assert.NoError(t, err)
			vars := map[string]string{"user_name": tt.args.userName}
			req = mux.SetURLVars(req, vars)
			req.Header.Add(helper.HeaderKeyContentType, "application/json")
			resp := httptest.NewRecorder()

			// test target
			ReportController(resp, req)
			assert.NoError(t, err)

			// assert db
			if tt.wantStatus == 200 {
				userReports, err := testingHelper.FindAllUserReports(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, 1, len(userReports))
				assert.Equal(t, userReports[0].ID, dummy.UserReport1.ID)
				assert.Equal(t, userReports[0].UserName, dummy.UserReport1.UserName)
				assert.Equal(t, userReports[0].Detail, dummy.UserReport1.Detail)
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

var successSubmitPostingReportReq = `
{
  "detail": "an inappropriate posting"
}
`
var errorSubmitPostingReportRespWithoutPostingID = `
{
  "status": 400,
  "message": "cannot be blank"
}
`
var errorSubmitPostingReportReqWithoutDetail = `
{
}
`
var errorSubmitPostingReportRespWithoutDetail = `
{
  "status": 400,
  "message": "detail: cannot be blank."
}
`
var errorSubmitPostingReportRespNotExistsData = `
{
  "status": 400,
  "message": "not exists data error"
}
`

func TestSubmitPostingReport(t *testing.T) {
	type args struct {
		postingID int
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
			args:       args{postingID: int(dummy.Posting1.ID), reqBody: successSubmitPostingReportReq},
			method:     http.MethodPost,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error path parameter without posting_id",
			args:       args{reqBody: successSubmitPostingReportReq},
			method:     http.MethodPost,
			want:       errorSubmitPostingReportRespWithoutPostingID,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body without detail",
			args:       args{postingID: int(dummy.Posting1.ID), reqBody: errorSubmitPostingReportReqWithoutDetail},
			method:     http.MethodPost,
			want:       errorSubmitPostingReportRespWithoutDetail,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error not exists data error",
			args:       args{postingID: int(dummy.Posting2.ID), reqBody: successSubmitPostingReportReq},
			method:     http.MethodPost,
			want:       errorSubmitPostingReportRespNotExistsData,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not allowed method",
			args:       args{postingID: int(dummy.Posting1.ID), reqBody: successSubmitPostingReportReq},
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
			userRepo := repository.NewUserRepository(db)
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)
			postingRepo := repository.NewPostingRepository(db)
			err = postingRepo.Create(context.Background(), &dummy.Posting1)
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/reports/postings/%s", strconv.Itoa(tt.args.postingID)), strings.NewReader(tt.args.reqBody))
			assert.NoError(t, err)
			vars := map[string]string{"posting_id": strconv.Itoa(tt.args.postingID)}
			req = mux.SetURLVars(req, vars)
			req.Header.Add(helper.HeaderKeyContentType, "application/json")
			resp := httptest.NewRecorder()

			// test target
			ReportController(resp, req)
			assert.NoError(t, err)

			// assert db
			if tt.wantStatus == 200 {
				postingReports, err := testingHelper.FindAllPostingReports(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, 1, len(postingReports))
				assert.Equal(t, postingReports[0].ID, dummy.PostingReport1.ID)
				assert.Equal(t, postingReports[0].PostingID, dummy.PostingReport1.PostingID)
				assert.Equal(t, postingReports[0].Detail, dummy.PostingReport1.Detail)
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
