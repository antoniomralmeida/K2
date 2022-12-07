@echo off
set GOARCH=amd64
set GOOS=linux
cd bin
go build  -o k2web.so ../k2web/k2web.go
go build  -o k2.so ../k2.go
cd ..
rem docker build . -f k2-back -t k2-back  --no-cache
rem docker build . -f k2-web -t k2-web  --no-cache
docker-compose build --no-cache
docker-compose up
