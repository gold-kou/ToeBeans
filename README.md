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

See backend/openapi.yml

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

- XSS 対策
  - トークンは Cookie に保存
- CSRF 対策
  - CSRF トークン
- Email 本人確認によるオプトイン
- パスワードはハッシュ化して保存
- HTTPS 対応

# Pricing memo

## GCP

- Cloud Vision API の　 LABEL_DETECTION
  - 1000 リクエスト/月まで無料
  - 500 万リクエスト/月まで 1.5$
  - https://cloud.google.com/vision/pricing

# Launch in local

## Backend
投稿画像が猫画像かどうかのチェックは本番環境のみでの使用を想定している。ローカルでも確認したい場合は、 `docker-compose.yml` の `APP_ENV` を `prd` にし、 `GOOGLE_API_KEY` の値を設定する。値は SSM のパラメータストアに設定されている。

```
$ cd backend
$ ./serverrun.sh
# go run main.go
```

`localhost:9000` で minio コンソールにログインし、 `toebeans-postings` の `Edit policy` で `* READ and Write` を追加する。サーバを再起動するたびに設定が必要。

## Frontend

`frontend/.env` をローカル用のものに変更する。

```
$ cd frontend
$ npm install
$ npm start
```

## Browser

Access to `localhost:3000/login` .
