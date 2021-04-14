# Development tips
## Launch servers in local
```
$ ./serverrun.sh
# go run main.go
```

If `listen tcp :8080: bind: address already in use exit status 1` happens, try below.

```
On app

# lsof -i:8080 -P
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
