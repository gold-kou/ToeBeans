package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	testingHelper "github.com/gold-kou/ToeBeans/backend/testing"
	"github.com/stretchr/testify/assert"
)

func TestGetHealthLiveness(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		want       string
		wantStatus int
	}{
		{
			name:       "success",
			method:     http.MethodGet,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "not allowed method",
			method:     http.MethodHead,
			want:       testingHelper.ErrNotAllowedMethod,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// http request
			req, err := http.NewRequest(tt.method, "/health/liveness", nil)
			assert.NoError(t, err)
			resp := httptest.NewRecorder()

			// test target
			HealthController(resp, req)
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

func TestGetHealthReadiness(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		want       string
		wantStatus int
	}{
		{
			name:       "success",
			method:     http.MethodGet,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "not allowed method",
			method:     http.MethodHead,
			want:       testingHelper.ErrNotAllowedMethod,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// http request
			req, err := http.NewRequest(tt.method, "/health/readiness", nil)
			assert.NoError(t, err)
			resp := httptest.NewRecorder()

			// test target
			HealthController(resp, req)
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
