@echo off
go mod tidy
set GOARCH=amd64
set GOOS=windows
rem docker run -d -p 9411:9411 --name zipkin openzipkin/zipkin-slim  
go build  -o ./bin/k2web.exe ./k2web/k2web.go
go build  -o ./bin/k2.exe k2.go
start .\bin\k2.exe
rem start .\bin\k2web.exe


