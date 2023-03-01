@echo off

set GOARCH=amd64
set GOOS=windows
set GCO_ENABLED=1
del bin\*.exe

set version="0.9.0-beta"
set build=%date:~6,4%-%date:~3,2%-%date:~0,2%-%time:~0,2%-%time:~3,2%-%time:~6,2%

git tag  %version% 
git push origin --tags

go build  -ldflags "-X 'github.com/antoniomralmeida/k2/pkg/version.Version=%version%' -X 'github.com/antoniomralmeida/k2/pkg/version.Build=%build%' " -o ./bin/k2web.exe ./cmd/k2web/main.go
go build  -ldflags "-X 'github.com/antoniomralmeida/k2/pkg/version.Version=%version%' -X 'github.com/antoniomralmeida/k2/pkg/version.Build=%build%' " -o ./bin/k2olivia.exe ./cmd/k2olivia/main.go
go build  -ldflags "-X 'github.com/antoniomralmeida/k2/pkg/version.Version=%version%' -X 'github.com/antoniomralmeida/k2/pkg/version.Build=%build%' " -o ./bin/k2.exe ./cmd/k2/main.go


del log\*.json
del log\*.log
start .\bin\k2olivia.exe
timeout 5
start .\bin\k2web.exe
.\bin\k2.exe
 timeout 10
taskkill /im  k2olivia.exe /f
taskkill /im  k2web.exe /f
