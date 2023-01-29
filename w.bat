@echo off
go get -u
go mod tidy
set GOARCH=amd64
set GOOS=windows
del bin\*.exe

set version="0.9.0-beta"
set build=%date:~6,4%-%date:~3,2%-%date:~0,2%-%time:~0,2%-%time:~3,2%-%time:~6,2%

git tag  %version% 
git push origin --tags

go build  -ldflags "-X 'github.com/antoniomralmeida/k2/version.version=%version%' -X 'github.com/antoniomralmeida/k2/version.build=%build%' " -o ./bin/k2web.exe ./k2web/k2web.go
go build  -ldflags "-X 'github.com/antoniomralmeida/k2/version.version=%version%' -X 'github.com/antoniomralmeida/k2/version.build=%build%' " -o ./bin/k2olivia.exe ./k2olivia/k2olivia.go
go build  -ldflags "-X 'github.com/antoniomralmeida/k2/version.version=%version%' -X 'github.com/antoniomralmeida/k2/version.build=%build%' " -o ./bin/k2.exe k2.go


del log\*.json
del log\*.log
start .\bin\k2olivia.exe
start .\bin\k2web.exe
.\bin\k2.exe
rem taskkill /im  k2olivia.exe /f
taskkill /im  k2web.exe /f
