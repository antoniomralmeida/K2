@echo off
go mod tidy
set GOARCH=amd64
set GOOS=windows
go build  -o ./bin/k2web.exe ./k2web/k2web.go
go build  -o ./bin/k2.exe k2.go
.\bin\k2.exe

