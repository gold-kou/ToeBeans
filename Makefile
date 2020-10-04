.PHONY: setup
setup:
	@go mod download
	@go mod verify

.PHONY: install
install:
	@go install github.com/gold-kou/ToeBeans/app/...

.PHONY: cross-install
cross-install:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -a -installsuffix cgo -ldflags '-w -extldflags "-static"' github.com/gold-kou/ToeBeans/app/...

.PHONY: openapi
openapi:
	@./hack/openapi/openapi-generate.sh

.PHONY: test
test:
	@go test -p=1 -covermode=count -coverprofile=cover.out github.com/gold-kou/ToeBeans/app/...
	@go tool cover -html=cover.out -o coverage.html

.PHONY: lint
lint:
	@./hack/lint.sh
