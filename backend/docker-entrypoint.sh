#!/bin/sh
# このファイルは削除するかも
# APP
export APP_ENV=$(/get-ssm-params -awsregion ap-​northeast-1 -extraparams /app_env | cut -d= -f2 | tr -d '"')
export DOMAIN=$(/get-ssm-params -awsregion ap-​northeast-1 -extraparams /domain | cut -d= -f2 | tr -d '"')
export LOG_LEVEL=$(/get-ssm-params -awsregion ap-​northeast-1 -extraparams /log_level | cut -d= -f2 | tr -d '"')
export JWT_SECRET_KEY=$(/get-ssm-params -awsregion ap-​northeast-1 -extraparams /jwt_secret_key | cut -d= -f2 | tr -d '"')
export CSRF_AUTH_KEY=$(/get-ssm-params -awsregion ap-​northeast-1 -extraparams /csrf_auth_key | cut -d= -f2 | tr -d '"')

# AWS
export AWS_ACCESS_KEY=$(/get-ssm-params -awsregion ap-​northeast-1 -extraparams /access_key | cut -d= -f2 | tr -d '"')
export AWS_SECRET_KEY=$(/get-ssm-params -awsregion ap-​northeast-1 -extraparams /secret_access_key | cut -d= -f2 | tr -d '"')

# DB
export DB_PASSWORD=$(/get-ssm-params -awsregion ap-​northeast-1 -extraparams /db_password | cut -d= -f2 | tr -d '"')

exec "$@"

