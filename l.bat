ren @echo off
go mod tidy
set GOARCH=amd64
set GOOS=linux
del bin\*.bin
go build  -o ./bin/k2web.bin ./k2web/k2web.go
go build  -o ./bin/k2olivia.bin ./k2olivia/k2olivia.go
go build  -o ./bin/k2.bin k2.go
del log\*.json
del log\*.log
cd docker
docker-compose build --no-cache
docker-compose up -d
cd ..