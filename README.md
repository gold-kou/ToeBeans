WIP

# ToeBeans
ToeBeans is a social media for cat-lovers!

# Tech stacks
## Frontend
HTML/CSS/JavaScript/React/SPA

## Backend
Go/OpenAPI/MySQL

## Infra
Docker/AWS(IAM/VPC/Route53/ECS/Fargate/RDS/S3/CloudFront/SES)/Terraform

## Tools
GitHub/GitHub Actions

# Infra

# Features
## Feature list
- User management
- Password reset
- Posting
- Like
- Comment
- Follow
- (Batch) Reset the count of sending email for password reset

## Not allowed features for guest user
- Like
- Posting
- Comment
- Follow
- Incident report
- Mute user
- Block user
- Direct message
- Other SNS sharing

Basically, guest user only can do read actions.

## Coming features
- Notification
- Other social media services sharing
- Refresh token
- Ranking
- Incident report
- Block user
- Movie posting

# Documents
## API
See openapi/openapi.yml

## RDB
https://docs.google.com/spreadsheets/d/1xIYH9PO4Hry3wTN6KYULvxmKMUQ6kwIWJNJJTyijZ_g/edit?usp=sharing

## Screen Layout
https://docs.google.com/presentation/d/1iqj8Hsm_CTQPWf_kTsZQMqlHoHxr-md7kc2Zsn8oom8/edit?usp=sharing

This is written in Japanese.

# Well designed points for performance
## Go
## RDB Indexing
## CloudFront
## Auto Scaling
## Pagenation

# Development tips
## Launch servers in local
```
$ ./serverrun.sh
# go run main.go
```

If the error of `listen tcp :8080: bind: address already in use exit status 1` happens, you might have failed to stop previous launching. 
Try below.

Kill process.

```
On app

# ps ax |grep main.go
# kill -9 <process>
```

Or another container might be block in local.
Remove unused container.

```
On Host

$ docker system prune -f
```

## Stop servers in local

```
# (control + Z)
^X^Z[1] + Stopped                    go run main.go

# ps ax |grep main.go

# kill -9 <process>

# exit
```

## curl


## UT
```
$ make test
```

Not enough test cases now.
