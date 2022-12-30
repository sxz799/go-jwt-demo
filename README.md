# go-jwt-study

```
go mod tidy
```


模拟登录

```
GET http://localhost:8000/login

{"username":"a","password":"a"}
```

登录后访问，access-token超时可自动更新有效期
refresh-token超时则提示需重新登录
```
GET http://localhost:8000/index

```

手动退出登录
```
GET http://localhost:8000/logout
```
