# Development tips
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

## generate OpenAPI models
```
$ make openapi
```

## Access log
example

```json
{"severity":"INFO","timestamp":"2021-06-13T12:06:35.760+0900","message":"","http_request":{"status":200,"method":"POST","host":"localhost:80","path":"/login","query":"","request_size":56,"remote_address":"172.22.0.1:35924","x_forwarded_for":"","user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36","referer":"http://localhost:3000/","protocol":"HTTP/1.1","latency":"20.6025ms"}}
```
