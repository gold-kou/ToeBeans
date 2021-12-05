# CORS
同一オリジンあるいは特別に許可しているオリジン（ローカル開発用）以外のオリジンからはリクエストを許可していません。

# CSRFトークン
「RFCで定められている安全」なメソッド（本サービスではGET）以外のリクエストは全て `X-CSRF-Token` ヘッダにCSRFトークンを格納する必要があります。秘密鍵を使って、CSRF

CSRFトークンは `/csrf-token GET` で取得可能です。

有効期限は24時間です。

[gorilla/csrf](https://github.com/gorilla/csrf) を使用しています。

# ID Token
ID Tokenは、本人であることを証明するトークンです。XSS対策として、Cookieに保存されてリクエストされることを想定しています。OpenIDConnectのID Tokenとは異なるものです。

トークンの生成と検証には、 [jwt-go](https://github.com/dgrijalva/jwt-go) を使用しています。

## トークンの中身
### header
HS256

### claim 
- iss
  - 発行者。文字列"ToeBeans"で固定。
- name
  - ユーザ名。
- iat
  - 発行時刻。
- exp
  - 有効期限。

## トークンの生成
秘密鍵を使って、上記内容からトークンを生成する。

## トークンの検証
1. Cookieの `id_token` プロパティに格納された値を取り出し、秘密鍵で検証できるかを確認する。
2. 有効期限を確認する。

# Development tips
## Login as a userA
email/password: userA@example.com/Password1234

## Launch servers in local
```
$ ./serverrun.sh
# go run main.go
```

If `listen tcp :80: bind: address already in use exit status 1` happens, try below.

```
On app

# lsof -i:80 -P
# kill -9 <process>
```

## UT
```
$ make test
```

## Generate OpenAPI models
```
$ make openapi
```

## Access log
example

```json
{"severity":"INFO","timestamp":"2021-06-13T12:06:35.760+0900","message":"","http_request":{"status":200,"method":"POST","host":"localhost:80","path":"/login","query":"","request_size":56,"remote_address":"172.22.0.1:35924","x_forwarded_for":"","user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36","referer":"http://localhost:3000/","protocol":"HTTP/1.1","latency":"20.6025ms"}}
```

## Be careful
It is impossible to request APIs in local by curl or tools because of `Forbidden - CSRF token invalid` .
Use frontend.