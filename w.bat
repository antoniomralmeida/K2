@echo off
go mod tidy
set GOARCH=amd64
set GOOS=windows
go build  -o ./bin/k2web.exe ./k2web/k2web.go
go build  -o ./bin/olivia.exe ./olivia/main.go
go build  -o ./bin/k2.exe k2.go
cd olivia
start ..\bin\olivia.exe
cd ..
start .\bin\k2.exe
start .\bin\k2web.exe
