package log

import (
	"bufio"
	"bytes"
	"net/http"
	"os"
	"text/template"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"gopkg.in/yaml.v3"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type accessLog struct {
	Method        string
	Host          string
	Path          string
	Query         string
	RequestSize   int64
	RemoteAddr    string
	XForwardedFor string
	Referer       string
	Protocol      string
}

type Logger struct {
	*zap.Logger
}

func (a *accessLog) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("method", a.Method)
	enc.AddString("host", a.Host)
	enc.AddString("path", a.Path)
	enc.AddString("query", a.Query)
	enc.AddInt64("request_size", a.RequestSize)
	enc.AddString("remote_address", a.RemoteAddr)
	enc.AddString("x_forwarded_for", a.XForwardedFor)
	enc.AddString("referer", a.Referer)
	enc.AddString("protocol", a.Protocol)
	return nil
}

func NewLogger() (*Logger, error) {
	tmpl, err := template.ParseFiles("/go/src/github.com/gold-kou/ToeBeans/backend/config/logger.yml.tpl")
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	err = tmpl.Execute(writer, os.Getenv)
	if err != nil {
		return nil, err
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	var config zap.Config
	if e := yaml.Unmarshal(buffer.Bytes(), &config); e != nil {
		return nil, e
	}
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{logger}, nil
}

func (l *Logger) LogHTTPAccess(r *http.Request) {
	accessLog := &accessLog{
		Method:        r.Method,
		Host:          r.Host,
		Path:          r.URL.Path,
		Query:         r.URL.RawQuery,
		RequestSize:   r.ContentLength,
		RemoteAddr:    r.RemoteAddr,
		XForwardedFor: r.Header.Get(helper.HeaderKeyXForwardedFor),
		Referer:       r.Referer(),
		Protocol:      r.Proto,
	}

	l.Info("", zap.Object("http_request", accessLog))
}
