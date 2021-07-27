# 構成図

# 構築手順
## 0. 前準備
クレデンシャル設定をする。

```
$ vi ~/.aws/credentials
```

## 1. ドメイン登録
`docs/operation_manual.md` の手順実施。

## 2. apply
```
$ cd environments/prd
$ AWS_PROFILE=tcpip-terraform GITHUB_TOKEN=xxx terraform apply -auto-approve
```

以下のエラーが発生した場合は、少し時間を置いてから再applyする。

```
Error: error creating ELBv2 Listener (arn:aws:elasticloadbalancing:ap-northeast-1:022111582403:loadbalancer/app/toebeans/460be5df191fb445): UnsupportedCertificate: The certificate 'arn:aws:acm:ap-northeast-1:022111582403:certificate/77af2cb4-5d0e-4359-812e-02119b9f32f7' must have a fully-qualified domain name, a supported signature, and a supported key size.
```

## 3. CodeStar設定
`docs/operation_manual.md` の手順を実施し、CodePipelineを再開する。

