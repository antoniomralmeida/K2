@echo off
set GOARCH=amd64
set GOOS=windows
cd bin
go build  -o k2web.exe ../k2web/k2web.go
go build  -o k2.exe ../k2.go
k2.exe
cd ..
