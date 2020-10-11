# this will be checked in CI
#!/usr/bin/env sh
set -Ceu

GENERATE_PATH='openapi/openapi-generated'

# code generate
docker run -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.1 generate \
  -i /local/openapi/openapi.yml \
  -g go-server \
  --output /local/${GENERATE_PATH}

# move
mkdir -p app/domain/model/http
sudo mv ${GENERATE_PATH}/go/model*.go app/domain/model/http

# rename package
sed -i -e 's/package openapi/package http/g' app/domain/model/http/model_*.go

# check openapi generated
git diff --exit-code app/domain/model/http

# check not tracked files
test -z "$(git ls-files --other --exclude-standard --directory --no-empty-directory)"
