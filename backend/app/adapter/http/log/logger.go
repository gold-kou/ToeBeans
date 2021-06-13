package log

import (
	"bufio"
	"bytes"
	"os"
	"text/template"
	"time"

	yaml "gopkg.in/yaml.v3"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AccessLog struct {
	Status        int
	Method        string
	Host          string
	Path          string
	Query         string
	RequestSize   int64
	RemoteAddr    string
	XForwardedFor string
	UserAgent     string
	Referer       string
	Protocol      string
	Latency       time.Duration
}

type Logger struct {
	*zap.Logger
}

func (a *AccessLog) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("status", a.Status)
	enc.AddString("method", a.Method)
	enc.AddString("host", a.Host)
	enc.AddString("path", a.Path)
	enc.AddString("query", a.Query)
	enc.AddInt64("request_size", a.RequestSize)
	enc.AddString("remote_address", a.RemoteAddr)
	enc.AddString("x_forwarded_for", a.XForwardedFor)
	enc.AddString("user_agent", a.UserAgent)
	enc.AddString("referer", a.Referer)
	enc.AddString("protocol", a.Protocol)
	enc.AddDuration("latency", a.Latency)
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

func (l *Logger) LogHTTPAccess(accessLog *AccessLog) {
	l.Info("", zap.Object("http_request", accessLog))
}
