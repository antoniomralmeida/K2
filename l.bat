@echo off
set GOARCH=amd64
set GOOS=linux
cd k2web
go build  -o ../bin/k2web.bin k2web.go 
cd ..
go build  -o ./bin/k2.bin k2.go 