@echo off
rem go get -u
rem go mod tidy
set GOARCH=amd64
set GOOS=windows
del bin\*.exe
<<<<<<< HEAD

set version="0.9.0-beta"
git tag -a %version% -m "version %version%"

go build  -ldflags "-X github.com/antoniomralmeida/k2/version.version=0.9.0-beta -X 'github.com/antoniomralmeida/k2/version.build=$(date)'" -o ./bin/k2web.exe ./k2web/k2web.go
go build  -ldflags "-X github.com/antoniomralmeida/k2/version.version=0.9.0-beta -X 'github.com/antoniomralmeida/k2/version.build=$(date)'" -o ./bin/k2olivia.exe ./k2olivia/k2olivia.go
go build  -ldflags "-X github.com/antoniomralmeida/k2/version.version=0.9.0-beta -X 'github.com/antoniomralmeida/k2/version.build=$(date)'" -o ./bin/k2.exe k2.go
=======
go build  -o ./bin/k2web.exe ./k2web/k2web.go
go build  -o ./bin/k2olivia.exe ./k2olivia/k2olivia.go
go build  -o ./bin/k2.exe k2.go
>>>>>>> 01887a253f097f28bcbfe9116bed04d1b593fab3
del log\*.json
del log\*.log
start .\bin\k2olivia.exe
start .\bin\k2web.exe
.\bin\k2.exe
rem taskkill /im  k2olivia.exe /f
rem taskkill /im  k2web.exe /f
