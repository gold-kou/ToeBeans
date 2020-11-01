#!/usr/bin/env bash

docker-compose -f docker-compose.test.yml run --rm app golangci-lint run
