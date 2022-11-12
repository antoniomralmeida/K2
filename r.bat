@echo off
set GOARCH=amd64
set GOOS=windows
cd k2web
go build  -o ../bin/k2web.exe k2web.go
cd ../bin
go run  ../k2.go
cd ..
