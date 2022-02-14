# local
FROM golang:1.15.0-buster as debugger
# install basic packages
RUN apt-get update && apt install -y --no-install-recommends \
  tzdata \
  wget \
  curl \
  lsof \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
  | sh -s v1.30.0 \
  && apt-get purge -y wget
# install go packages
WORKDIR /realize
RUN go mod init vendor
RUN echo 'replace gopkg.in/urfave/cli.v2 => github.com/urfave/cli/v2 v2.1.1' >> go.mod
RUN go get github.com/oxequa/realize \
  && go get github.com/rakyll/gotest \
  && go get github.com/jwilder/dockerize
WORKDIR /go/src/github.com/gold-kou/ToeBeans/backend
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY go.mod go.sum ./
RUN GO111MODULE=on go mod download
COPY . .
# launch
EXPOSE 80
CMD ["dockerize", "-wait", "tcp://db:3306", "-timeout", "60s", "realize", "start", "--run", "--no-config"]

# 本番で使うバイナリのビルド用環境
FROM golang:1.15.0-buster as builder
RUN apt update && apt install -y --no-install-recommends git tzdata
WORKDIR /go/src/github.com/gold-kou/ToeBeans/backend
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY go.mod go.sum ./
RUN GO111MODULE=on go mod download
COPY . .
RUN go build

# 本番の実行環境
FROM debian:10.9-slim as executor
RUN apt update && apt upgrade -y --no-install-recommends liblz4-1 \
  && rm -rf /var/lib/apt/lists/*
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /go/src/github.com/gold-kou/ToeBeans/backend
COPY config/logger.yml.tpl config/logger.yml.tpl
WORKDIR /
COPY --from=builder /go/src/github.com/gold-kou/ToeBeans/backend/backend /backend
# 基本的に値を変えない環境変数のみをここに設定
ENV AWS_REGION=ap-​northeast-1 DB_NAME=toebeansdb DB_PORT=3306 DB_USER=toebeans SYSTEM_EMAIL=no-reply@toebeans.tk TZ=Asia/Tokyo
EXPOSE 80
CMD ["/backend"]