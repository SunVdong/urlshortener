# 长连接转短链接

这是一个后端接口服务

## 依赖于
- 开发工具
```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```
- 数据库： postgres
- 缓存： redis
  
## 运行
```
make migrate_up
sqlc generate
go mod tidy
go run main.go
```
