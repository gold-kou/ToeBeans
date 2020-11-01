#!/bin/sh

# DB
export MYSQL_DATABASE=
export MYSQL_USER=
export MYSQL_PASSWORD=$(/get-ssm-params -awsregion ap-southeast-1 -extraparams  | cut -d= -f2 | tr -d '"')
export MYSQL_HOST=$(/get-ssm-params -awsregion ap-southeast-1 -extraparams  | cut -d= -f2 | tr -d '"')


# AWS


exec "$@"

