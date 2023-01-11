@echo off
go mod tidy
set GOARCH=amd64
set GOOS=linux
go build  -o ./bin/k2web.so ./k2web/k2web.go
go build  -o ./bin/k2olivia.so ./k2olivia/k2olivia.go
go build  -o ./bin/k2.so k2.go
del log\*.json
del log\*.log
docker-compose build --no-cache
docker-compose up -d
