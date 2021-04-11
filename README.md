WIP

# ToeBeans
ToeBeans is a social media service for cat-lovers!

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

Basically, guest user only can do read actions.

## Coming features
- Notification
- Other SNS sharing
- Refresh token
- Ranking
- Incident report
- Direct message
- Block user
- Movie posting

# Documents
## API
See backend/openapi/openapi.yml

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

# セキュリティ面で工夫しているところ
- XSS対策
  - トークンはCookieに保存
- CSRF対策
  - CSRFトークン
- オプトイン
  - Email確認
- パスワードはハッシュ化して保存
- HTTPS対応