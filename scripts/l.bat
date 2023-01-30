ren @echo off
go mod tidy
set GOARCH=amd64
set GOOS=linux
del bin\*.bin
set version="0.9.0-beta"
set build=%date:~6,4%-%date:~3,2%-%date:~0,2%-%time:~0,2%-%time:~3,2%-%time:~6,2%

git tag  %version% 
git push origin --tags

go build  -ldflags "-X github.com/antoniomralmeida/k2/version.version=%version% -X github.com/antoniomralmeida/k2/version.build=%build% " -o ./bin/k2web.bin ./k2web/k2web.go
go build  -ldflags "-X github.com/antoniomralmeida/k2/version.version=%version% -X github.com/antoniomralmeida/k2/version.build=%build% " -o ./bin/k2olivia.bin ./k2olivia/k2olivia.go
go build  -ldflags "-X github.com/antoniomralmeida/k2/version.version=%version% -X github.com/antoniomralmeida/k2/version.build=%build% " -o ./bin/k2.bin k2.go

del log\*.json
del log\*.log
cd docker
docker-compose build --no-cache
docker-compose up -d
cd ..