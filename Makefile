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
	@docker-compose -f docker-compose.test.yml build app
	@docker-compose -f docker-compose.test.yml run --rm app dockerize -wait tcp://db-test:3306 -timeout 60s gotest -p 1 -v github.com/gold-kou/ToeBeans/app/...
	@docker-compose -f docker-compose.test.yml down --remove-orphans

.PHONY: lint
lint:
	@./hack/lint.sh
