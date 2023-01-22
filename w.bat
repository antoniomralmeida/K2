@echo off
go mod tidy
set GOARCH=amd64
set GOOS=windows
del bin/*.exe
go build  -o ./bin/k2web.exe ./k2web/k2web.go
go build  -o ./bin/k2olivia.exe ./k2olivia/k2olivia.go
go build  -o ./bin/k2.exe k2.go
del log\*.json
del log\*.log
start .\bin\k2olivia.exe
start .\bin\k2web.exe
.\bin\k2.exe
taskkill /im  k2olivia.exe /f
taskkill /im  k2web.exe /f
