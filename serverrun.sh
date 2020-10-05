#!/bin/bash
base=${0##*/}

docker-compose up -d --build
docker-compose exec app sh
docker-compose down