# ToeBeans
ToeBeans is a social media service for cat-lovers!

TODO demo movie

# Features
## Released
Regular user can:

- User
  - Register a user
  - Read user info
  - Update user info
  - Delete a user
  - Password change
- Posting
  - Post a posting
  - Read postings
  - Delete a posting
- Like
  - Do a like
  - Delete a like
- Follow
  - Do a follow
  - Delete a follow
- Report
  - Do a report

Guest user can:

- User
  - Read user info
- Posting
  - Post a posting
  - Read postings
  - Delete a posting
- Like
  - Do a like
  - Delete a like
- Follow
  - Do a follow
  - Delete a follow

Basically, guest user can do only read actions.

## Coming
- Comment
- Search
- Notification
- Refresh token
- Ranking
- Block user
- Movie posting

# Tech stacks
- HTML/CSS/JavaScript/React
- Go/OpenAPI/MySQL
- Docker/AWS/Terraform
- GitHub/GitHub Actions

# Architecture Layout
![Architecture](material/ToeBeans%20Architecture.drawio.png)

※インフラコスト削減を優先しているためRDSの冗長化はしていません。

# Documents
## API
See `backend/openapi.yml` .

## Backend
See `backend/README.md` .

## Infra
See `infra/README.md` .

## Design memo
https://docs.google.com/presentation/d/1iqj8Hsm_CTQPWf_kTsZQMqlHoHxr-md7kc2Zsn8oom8/edit?usp=sharing

# Launch in local
## Backend
投稿画像が猫画像かどうかのチェックは本番環境のみでの使用を想定している。ローカルでも確認したい場合は、`docker-compose.yml`の`APP_ENV`を`prd`にし、`GOOGLE_API_KEY`の値を設定する。値はSSMのパラメータストアに設定されている。

```
$ cd backend
$ ./serverrun.sh
# go run main.go
```

`localhost:9000`でminioコンソールにログインし、`toebeans-postings`の`Edit policy`で`* READ and Write`を追加する。サーバを再起動するたびにこの設定は必要。

## Frontend
`frontend/.env`をローカル用のものに変更する。

```
$ cd frontend
$ npm install
$ npm start
```

## Browser
Access to `localhost:3000/login`.
