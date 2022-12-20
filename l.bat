@echo off
set GOARCH=amd64
set GOOS=linux
cd bin
go build  -o k2web.so ../k2web/k2web.go
go build  -o k2.so ../k2.go
cd ..
docker-compose build --no-cache
docker-compose up -d
