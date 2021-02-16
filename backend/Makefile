.PHONY: setup
setup:
	@go mod download
	@go mod verify

.PHONY: install
install:
	@go install github.com/gold-kou/ToeBeans/backend/app/...

.PHONY: openapi
openapi:
	@./hack/openapi/openapi-generate.sh

.PHONY: test
test:
	@docker-compose -f docker-compose.test.yml build app
	-@docker-compose -f docker-compose.test.yml run --rm app dockerize -wait tcp://db-test:3306 -timeout 60s dockerize -wait tcp://minio:9000 -timeout 60s gotest -p 1 -v github.com/gold-kou/ToeBeans/backend/app/...
	@docker-compose -f docker-compose.test.yml down --remove-orphans

.PHONY: lint
lint:
	@docker-compose -f docker-compose.test.yml build app
	-@docker-compose -f docker-compose.test.yml run --rm app golangci-lint run --config /go/src/github.com/gold-kou/ToeBeans/backend/.golangci.yml -v
	@docker-compose -f docker-compose.test.yml down --remove-orphans
