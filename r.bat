@echo off
cd k2web
go build   k2web.go
cd ..
go run  k2.go