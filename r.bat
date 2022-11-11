@echo off
cd k2web
go build -o mybinary -ldflags "-X main.version=0.5.0 -X 'main.build=$(date)'" k2web.go

cd ..
go build -o mybinary -ldflags "-X main.version=0.5.0 -X 'main.build=$(date)'"  k2.go

go run -ldflags "-X main.version=0.5.0 -X main.build=$NT_GNU_BUILD_ID" k2.go