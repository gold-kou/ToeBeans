# Architecture
ヘキサゴナルアーキテクチャを採用しています。

ヘキサゴナルアーキテクチャは、ドメイン領域を中心に見据えてそのほかを外側に押しやることで、コンポーネント間を疎結合にし、以下のメリットをもたらします。

- 変更に強くなる
  - 例えば、DBがMySQLからPostgreSQLへ変わったとしても、Adapterの実装とAdapterの呼び出し部分を書き換えるだけで済みます。つまり、Application側はアダプタを付け加えるだけで済みます。
- テストを書きやすい
  - Adapterをモックに置き換えることで外部と接続する機能のテストも容易に行えます。

# Security
## Authorization
IDトークンによる認証方式としています。IDトークンは、本人であることを証明するトークンです。JWTの仕様に従っています。認証を必要とするAPIでは、Cookieに正しいIDトークンが格納されていない場合は401エラーを返します。

OpenIDConnectのIDトークンとは異なるものです。

### IDトークンの中身
#### Header
署名検証を行うための以下のメタ情報をBase64エンコードしたものです。

- typ
  - JWT
- alg
  - 署名アルゴリズム
    - HS256

#### Payload(Claims)
IDトークンの中身となる以下の情報をBase64エンコードしたものです。

- iss
  - 発行者名。文字列"ToeBeans"。
- sub
  - ユーザID
- name
  - ユーザ名
- iat
  - 発行時刻
- exp
  - 有効期限

#### Signature
`Header.Payload` をalgの署名アルゴリズムで署名し、Base64エンコードしたものです。

### IDトークンの生成と検証
[jwt-go](https://github.com/dgrijalva/jwt-go)を使用しています。

#### 生成
1. `jwt.New` でHeaderを含めた初期化を行います。typがJWTであることは自明なので署名アルゴリズム（HS256）だけを指定します。
2. Claimsを設定します。
3. 署名に必要な鍵をjwt-goの関数に渡してトークンを生成します。

```go
func GenerateToken(userID int64, userName string) (tokenString string, err error) {
	// header
	token := jwt.New(jwt.SigningMethodHS256)

	// claims
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "ToeBeans"
	claims["sub"] = strconv.Itoa(int(userID))
	claims["name"] = userName
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * TokenExpirationHour).Unix()

	// generate token by secret key
	tokenString, err = token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
```

#### 検証
1. 署名アルゴリズムが一致しているかの確認します。
2. 秘密鍵で復号できるかを確認します。
3. 有効期限を確認します。

```go
func VerifyToken(tokenString string) (userID int64, userName string, err error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err = errUnexpectedSigningMethod
			return nil, err
		}
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				err = errTokenExpired
				return
			}
			err = errTokenInvalid
			return
		}
	}
  // ...
}
```

## XSS対策
IDトークンはLocalStorageでなく、Cookie(httpOnly)に格納し、JavaScriptから参照・更新ができないようにしています。Cookie(httpOnly)への格納は完全にXSSを対策できるものではありませんが、簡単に実装できるため、緩和策として講じています。

## CSRF対策
### 背景
IDトークンをCookieに格納しているため、CSRF攻撃によりユーザに身に覚えないのリクエストが実行されてしまう可能性があります。

### CSRFトークン
CSRFトークンによるCSRF対策をしています。RFCで定められている「安全」なメソッド（本サービスではGET）以外のリクエストは全て `X-CSRF-Token` ヘッダにCSRFトークンが格納されている必要があります。

CSRFトークンは `/csrf-token GET` で取得可能ですが、CORSで許可されているオリジン（SPAのフロントエンド）からしか取得できません。

有効期限は24時間です。

### CSRFトークンの生成と検証
[gorilla/csrf](https://github.com/gorilla/csrf) を使用しています。

#### 生成
Token関数でCSRFトークンを生成しています。ライブラリ内で乱数を使って生成されています。これをAPIレスポンスで返します。

```go
token := csrf.Token(r)
```

#### 検証
Protect関数でCSRFトークンを検証しています。検証は非安全なHTTPメソッドの場合のみ、ミドルウェアで実施します。鍵は環境変数で与えたものを使用しています。

```go
csrfMiddleware := csrf.Protect([]byte(csrfAuthKey))
```

## CORS
JavaScriptのSame Origin Policyにより、異なるオリジンからのリクエストは許可されていません。しかしながら、SPA＋APIサーバの構成のため、フロントエンドとバックエンドがそれぞれ異なるドメイン上に存在します。そこで、CORSにより `toebeans.ml` と `localhost:3000` からのリクエストは `Access-Control-Allow-Origin` を設定し、許可するようにしています。加えて、IDトークンをCookieに格納しているため、異なるオリジンへのCookieを許可するために `Access-Control-Allow-Credentials: true` にしています。

## メール本人確認によるダブルオプトイン
アカウント作成時に登録されたメールアドレス宛に、アカウントを有効化するためのアクティベーションキーが付与されたURLリンクを記載したメールを送信しています。そのリンクが踏まれることで、アカウント作成が完了します。ダブルオプトインにより、他人のメールアドレスが使用されたり、存在しないメールアドレスが使用されることを防いでいます。

## パスワードのハッシュ化
パスワードはハッシュ化したうえでRDBに保存しています。

ハッシュ化には[crypt/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)パッケージのGenerateFromPassword関数を使用しています。

```go
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.reqRegisterUser.Password), bcrypt.DefaultCost)
```

パスワード照合時には、DBに保存されたハッシュ化済みパスワードとログイン時に入力されたパスワードを照合します。同パッケージのCompareHashAndPassword関数を使用しています。

```go
if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(l.reqLogin.Password)); err != nil {
	return "", ErrNotCorrectPassword
}
```

# Development tips
## UT
```
$ make test
```

## Launch servers in local
### (docker-composeを使う場合)コンテナとアプリケーション起動
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

### (ECRのイメージを使う場合)コンテナとアプリケーション起動
```
$ $(aws ecr get-login --profile tcpip-terraform --region ap-northeast-1 --no-include-email)
WARNING! Using --password via the CLI is insecure. Use --password-stdin.
Login Succeeded
$ docker pull XXX.dkr.ecr.ap-northeast-1.amazonaws.com/toebeans:latest
$ docker image ls |grep toebeans
$ docker run -it b83a48f97efb /bin/bash
```

### Table migration
`docker-compose.yml` および `docker-compose.test.yml` の volumes で `./toebeans-sql/mysql/entrypoint:/docker-entrypoint-initdb.d` としているため、コンテナ起動時に自動でマイグレーションが実行される。ただし、前回のが残っている場合は、volumeは削除する必要がある。

```
$ docker volume ls |grep backend_db
local     backend_db
local     backend_db-test

$ docker volume rm backend_db
backend_db

$ docker volume rm backend_db-test
backend_db-test
```

### Login as a userA
`userA` is automatically created by `toebeans-sql/mysql/entrypoint/002_insert_dummy_data.sql`.

email/password: userA@example.com/Password1234

## Access log
```json
{
  "severity": "INFO",
  "timestamp": "2021-06-13T12:06:35.760+0900",
  "message": "",
  "http_request": {
    "status": 200,
    "method": "POST",
    "host": "localhost:80",
    "path": "/login",
    "query": "",
    "request_size": 56,
    "remote_address": "172.22.0.1:35924",
    "x_forwarded_for": "",
    "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36",
    "referer": "http://localhost:3000/",
    "protocol": "HTTP/1.1",
    "latency": "20.6025ms"
  }
}
```
