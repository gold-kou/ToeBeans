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

## UT
```
$ make test
```
