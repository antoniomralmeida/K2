@echo off
go mod tidy
set GOARCH=amd64
set GOOS=linux
del bin\*.bin
set version="0.9.0-beta"
set build=%date:~6,4%-%date:~3,2%-%date:~0,2%-%time:~0,2%-%time:~3,2%-%time:~6,2%

git tag  %version% 
git push origin --tags

go build  -ldflags "-X github.com/antoniomralmeida/k2/version.version=%version% -X github.com/antoniomralmeida/k2/version.build=%build% " -o ./bin/k2web.bin ./cmd/k2web/main.go
go build  -ldflags "-X github.com/antoniomralmeida/k2/version.version=%version% -X github.com/antoniomralmeida/k2/version.build=%build% " -o ./bin/k2olivia.bin ./cmd/k2olivia/main.go
go build  -ldflags "-X github.com/antoniomralmeida/k2/version.version=%version% -X github.com/antoniomralmeida/k2/version.build=%build% " -o ./bin/k2.bin ./cmd/k2/main.go

del log\*.json
del log\*.log
docker-compose -f ./build/docker-compose.yml  build
docker-compose -f ./build/docker-compose.yml  up -d  
