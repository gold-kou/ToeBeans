# local
FROM golang:1.15.0-buster as debugger
# install basic packages
RUN apt-get update && apt install -y --no-install-recommends \
    tzdata \
    wget \
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
WORKDIR /go/src/github.com/gold-kou/ToeBeans
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY go.mod go.sum ./
RUN GO111MODULE=on go mod download
COPY . .
# launch
EXPOSE 8080
CMD ["dockerize", "-wait", "tcp://db:3306", "-timeout", "60s", "realize", "start", "--run", "--no-config"]
