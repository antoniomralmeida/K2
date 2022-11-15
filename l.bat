@echo off
set GOARCH=amd64
set GOOS=linux
cd bin
go build  -o k2web.bin ../k2web/k2web.go
go build  -o k2.bin ../k2.go
cd ..
docker build . -t k2-app-back  --no-cache
docker-compose up
