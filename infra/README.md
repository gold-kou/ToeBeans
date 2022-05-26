# 構築手順
## ドメイン登録
### freenomでドメイン取得
適当に画面入力してドメインを取得する。

無料枠最大の12ヶ月で取得する。更新も無料。ただし更新を忘れるとそのドメインは有料になってしまうのでカレンダーに登録するなどしよう。更新メールも最近は送られなくなった？

`Forward this domain or Use DNS` の選択は特にせずにContinueしてしまってよい。

### Route53でホストゾーンの作成
`ホストゾーンの作成` に進み、freenom で取得したドメイン名を入力する。

![](docs/images/Route_53_ドメイン作成.png)

### freenom でレコード設定
`Services > My Domains > Manage Domain > Manage Tools > Use custom nameservers (enter below)` と進む。

Route53で払い出された該当ドメインNSレコードの値4つを入力する。
-> `Change Nameservers` を押す。登録後、自動で大文字になる。

### ref
- https://note.com/dafujii/n/n406f385651e2
- https://note.com/dafujii/n/n12bb564081f1

## メールサーバ構築
### 新コンソール手順
#### ドメイン検証からテストメール送信まで
https://www.blog.danishi.net/2022/02/13/post-5807/

- Verified identitiesへ進む。
- ドメイン検証
  - Create identityをクリック。Domainを選択。Assign a default configuration setを選択。Create identityをクリック。
    - CNAMEのRoute53登録は自動でされる。
- メール検証
  - Create identityをクリック。Email addressを選択。メールアドレスを入力。Create identityをクリック。
- テストメール送信
  - 対象ドメインを選択。Send test emailをクリック。From-addressは@前を適当に入力。ScenarioでCustomを選択し、送信先メールアドレスとして↑で検証したメールアドレスを入力。SubjectとBodyを適当に入力してSend test emailをクリック。

#### Sandboxから抜ける
https://docs.aws.amazon.com/ja_jp/ses/latest/dg/request-production-access.html

※Classic Console時に既に設定していたため、実際には↑の手順は試していない。

### Classic Console手順（廃止）
※Classic Consoleが廃止されたため以下手順は使えない。

#### ref
- https://note.com/dafujii/n/n0365dc0a89af
- https://docs.aws.amazon.com/ses/latest/DeveloperGuide/request-production-access.html

#### ドメイン検証とレコード設定
SES コンソールの `Domains` の `verify a New Domain` を押す。
freenom で取得したドメイン `toebeans.ml` を入力して、 `Generate DKIM Settings` にチェックして、 `Verify This Domain` を押す。

![](docs/images/verify_new_domain.png)

ポップアップが表示されるので、 `Use Route53` を押す。

`Email Receiving Record`にチェックを入れ、`Create Record Sets` ボタンを押す。

![](docs/images/SES_レコード設定.png)

#### ドメイン有効確認
有効になっていることを確認する。72 時間以内に有効となるらしいが、たいていは数分以内で完了する。ただし、DKIM だけ一時間ほどかかるかも。

`Email Addresses > Verify a New Email Address` で適当なメールアドレスを入力し、検証する。

![](docs/images/verify_new_email.png)

メールのテスト送信をする。さきほど検証したメールアドレスを送信先に指定する。

![](docs/images/send_test_email.png)

#### sandbox から抜ける
SES は作成時点では sandbox のまま。
このままでは検証されたメールに対してのみにしか送信できないため、sandbox から抜けるための申請が必要。

`Email Sending > Sending Statistics > Edit your account details` と進む。

`Enable production access` を Yes にして申請内容を適当に埋める。今回は、ユーザ登録機能にメール送信が必要なので、Sandox から抜けたいです的なことを拙い英語で書いた。

![](docs/images/SES_sandboxから抜ける申請.png)

## Terraform
### 前提
Terraformインストール済み。

### IAM ユーザ作成
AWS コンソールにてTerraform実行用のユーザをコンソールで作成する。

### クレデンシャル設定
```
$ vi ~/.aws/credentials
[terraform]
aws_default_region=ap-northeast-1
aws_access_key_id=XXXXX
aws_secret_access_key=XXXXX
```

### tfstate管理用バケットの作成
AWS コンソールにてS3バケットを作成する。「パブリックアクセスをすべてブロック」・「バージョニングを有効」・「暗号化を有効」で作成する。バケット名は、 `main.tf` の bucket と一致させる。

### GitHubとAWSの連携設定
`AWS Connector for GitHub` を設定する。（手順詳細忘れた）。

### SSHキーペアの作成
`ssh-keygen -t rsa -f toebeans -N ''` とかだった気がする。

### GitHub トークンの生成
GitHub コンソールの `Settings > Developer Settings > Personal access tokens` と進む。
新規の場合は、 `Generate new token` を押す。 `Select scopes` では `repo` と `admin:repo_hook` を全てチェック。
再発行の場合は、 `Regenarate token` を押す。

### SSM パラメータストアの設定
`cotainer_definitions.json` の `secrets` の内容を設定する。

`github_token` に関しては apply 前に設定が必須。

`google_api_key` の値は `$ cat backend/secret/service-account.json | tr -d '\n'` の実行結果からスペースを全て削除したものを設定する。
`db_host` の値は RDS エンドポイントの値を設定する。

### S3 バケット名の決定
Goアプリケーションが使用するS3バケット名を決定し、 `s3.tf` を編集する。バケット名は世界中で一意である必要がある。諸事情により、複数のAWSアカウントで開発しているため、バケット名が重複してしまう場合はバケット名を既存のものから変更すること。その際は、 `backend/Dockerfile` のENVの設定もあわせて実施すること。

### apply
```
$ cd environments/prd
$ AWS_PROFILE=terraform GITHUB_TOKEN=xxx terraform apply -auto-approve
```

以下のエラーが発生した場合は、少し時間を置いてから再 apply する。
// TODO depends_on 使えば回避できるかもしれない。

```
Error: error creating ELBv2 Listener (arn:aws:elasticloadbalancing:ap-northeast-1:022111582403:loadbalancer/app/toebeans/460be5df191fb445): UnsupportedCertificate: The certificate 'arn:aws:acm:ap-northeast-1:022111582403:certificate/77af2cb4-5d0e-4359-812e-02119b9f32f7' must have a fully-qualified domain name, a supported signature, and a supported key size.
```

## CodeStar接続
初回apply時はCodePipelineの実行に失敗してしまうため、以下の手順で再開する。

接続の項目から対象のコネクションを選択する。
![](docs/images/AWS_Developer_Tools_1.png)

`保留中の接続を更新` を押す。
![](docs/images/AWS_Developer_Tools_2.png)

GitHubアプリで自分のアカウントを選択し、 `接続` を押す。
![](docs/images/AWS_Developer_Tools_3.png)

対象のパイプラインを選択し、 `変更をリリースする` を押す。

## DBの変更
adminユーザのパスワード変更、テーブルマイグレーション、アプリケーション用ユーザの作成を実施します。

1. EC2 コンソールを利用し、踏み台サーバにログインする。
2. `mysql -u admin -p -h <RDSエンドポイント>` を実行する。RDSエンドポイントはコンソールの `接続とセキュリティ` から確認可能。初期パスワードはvariables.tfを参照する。
3. `SET PASSWORD = PASSWORD('XXXXX');` を実行してadminユーザのパスワードを変更する。パスワード値は任意の値。
4. `CREATE DATABASE toebeansdb DEFAULT CHARACTER SET utf8;` を実行する。
5. `USE toebeansdb;` を実行する。
6. `backend/toebeans-sql/mysql/entrypoint/001_create_tables.sql` の内容を実行する。
7. `backend/toebeans-sql/mysql/create_user.sql` の内容を実行する。パスワード値は任意の値。
8. `backend/toebeans-sql/mysql/entrypoint/002_insert_dummy_data.sql` の内容を実行する。ゲストユーザのみでよい。
9. exitする。

## CloudFront修正
`cycle error` によりACMをTerraformのコード上で指定できない都合上、コンソールで設定の追加をする必要がある。設定後数分で403Errorでなくなる。

- 代替ドメイン名を追加してtoebeans.mlを入力する
- カスタムSSL証明書でバージニアのものを選択する

![](docs/images/CloudFront_add.png)

※Terraform のコード上で `cloudfront_default_certificate = true` としているめ、applyされるたびに本設定をやり直す必要がある。

# Deploy
## Frontend
- `frontend/` に移動。
- `npm run build` を実行する。
- `build` ディレクトリ配下のファイルとディレクトリをS3へアップロードする。
  - build ディレクトリ自体はアップロードしないように注意する。

## Backend
`gold-kou/toebeans` のmasterブランチにマージされると、CodePipelineのSourceとBuildによりECRへ最新のバックエンドのDockerイメージがpushされ、Deployにより自動でデプロイされる。マージしてからデプロイされるまでには約20分かかる。

# Pricing memo
## GCP
- Cloud Vision APIのLABEL_DETECTION
  - 1000リクエスト/月まで無料
  - 500万リクエスト/月まで1.5$
  - https://cloud.google.com/vision/pricing

## AWS
リクエスト量にもよるが、起動コストだけでおそらく2万円程度。

個人開発にしては料金コストが重いなと感じたことをメモします。

- AWS Config
  - ECS上で起動失敗するアプリケーションをずっとオートヒーリングし続けてしまうことで料金がかかった経験がある。
    - Essential container in task exited でタスクが起動せず、desired_count=2にしていたため。
- RDS
  - 起動料金が高い。
  - 開発中はこまめに落とすのも手。
    - 起動に20分くらいかかるので効率悪いかも。
    - 7日間停止していると自動で8日目に起動するので要注意。
- NAT Gateway
  - 起動料金が高い。
  - プライベートインスタンスがインターネットに接続するためにサービス運用中は必要だが、開発中はこまめに落とす。
- VPC Endpoint
  - 起動料金が高い。
    - NAT Gatewayの費用削減になるかもと思い、使ってみたがあまりうまくいかなかったのでやめた。
- ELB
  - 起動料金が高い。
  - こまめに落とす。

開発中は節約のために以下を実施している。

- RDSインスタンスの停止
  - 毎週忘れないようにする
- ECSのdesired_countを0にしてapply
- NAT Gateway関連のネットワークリソースをコメントアウトしてapply
