package lib

import (
	"os"
	"strings"
	"testing"

	"github.com/gold-kou/ToeBeans/testing/dummy"

	"github.com/stretchr/testify/assert"
)

var sharedTestToken string

func TestGenerateToken(t *testing.T) {
	type args struct {
		userName string
	}
	tests := []struct {
		name        string
		args        args
		environment string
		want        string
		wantErr     bool
	}{
		{
			name:        "success",
			args:        args{userName: dummy.User1.Name},
			environment: dummy.SecretKey,
			want:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			// set env
			tmp := setTestEnv("JWT_SECRET_KEY", tt.environment)
			defer tmp()

			// test target
			got, err := GenerateToken(tt.args.userName)
			sharedTestToken = got

			if tt.wantErr {
				a.Error(err)
			} else {
				a.NoError(err)
				// just checking HEADER
				a.Equal(tt.want, strings.Split(got, ".")[0])
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	type args struct {
		tokenString string
	}
	tests := []struct {
		name        string
		args        args
		environment string
		want        string
		wantErr     bool
		watnErrMsg  string
	}{
		{
			name:        "success",
			args:        args{tokenString: sharedTestToken},
			environment: dummy.SecretKey,
			want:        dummy.User1.Name,
			wantErr:     false,
		},
		{
			name:        "fail(expired)",
			args:        args{tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM3MjIzNjIsImlhdCI6IjIwMjAtMDMtMDhUMTE6NTI6NDIuMjIxMjY2NCswOTowMCIsIm5hbWUiOiJ0ZXN0In0.FYMJmXo17aUhTpdaLifMovDQ0BiKSq8LnssLwxFvshI"},
			environment: dummy.SecretKey,
			want:        dummy.User1.Name,
			wantErr:     true,
			watnErrMsg:  "token is expired",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set env
			tmp := setTestEnv("JWT_SECRET_KEY", tt.environment)
			defer tmp()

			a := assert.New(t)

			// test target
			got, err := VerifyToken(tt.args.tokenString)

			// assert
			if tt.wantErr {
				a.Error(err)
				a.EqualError(err, tt.watnErrMsg)
			} else {
				a.NoError(err)
				a.Equal(tt.want, got)
			}
		})
	}
}

func setTestEnv(key, val string) func() {
	preVal := os.Getenv(key)
	os.Setenv(key, val)
	return func() {
		os.Setenv(key, preVal)
	}
}
