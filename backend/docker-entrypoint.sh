#!/bin/sh
# APP
export APP_ENV=$(/get-ssm-params -awsregion ap-​northeast-1 -extraparams /app_env | cut -d= -f2 | tr -d '"')
export DOMAIN=toebeans.tk
export TZ=Asia/Tokyo
export LOG_LEVEL=warn
export JWT_SECRET_KEY=$(/get-ssm-params -awsregion ap-​northeast-1 -extraparams /jwt_secret_key | cut -d= -f2 | tr -d '"')
export CSRF_AUTH_KEY=$(/get-ssm-params -awsregion ap-​northeast-1 -extraparams /csrf_auth_key | cut -d= -f2 | tr -d '"')

# GCP
export GOOGLE_APPLICATION_CREDENTIALS=secret/service-account.json

# AWS
export AWS_ACCESS_KEY=$(/get-ssm-params -awsregion ap-​northeast-1 -extraparams /access_key | cut -d= -f2 | tr -d '"')
export AWS_SECRET_KEY=$(/get-ssm-params -awsregion ap-​northeast-1 -extraparams /secret_access_key | cut -d= -f2 | tr -d '"')
export AWS_REGION=ap-​northeast-1
export S3_BUCKET_POSTINGS=/toebeans-postings
export S3_BUCKET_ICONS=/toebeans-icons
export SYSTEM_EMAIL=no-reply@toebeans.tk

# DB
export DB_NAME=toebeansdb
export DB_USER=toebeans
export DB_PASS=$(/get-ssm-params -awsregion ap-​northeast-1 -extraparams /db_pass | cut -d= -f2 | tr -d '"')
export DB_HOST=db
export DB_PORT=3306

exec "$@"

