@echo off
set GOARCH=amd64
set GOOS=linux
cd bin
go build  -o k2web.so ../k2web/k2web.go
go build  -o k2.so ../k2.go
cd ..
scp -i "..\..\manoel.pem" config/*  ec2-user@ec2-15-228-21-212.sa-east-1.compute.amazonaws.com:/home/ec2-user/config 
scp -i "..\..\manoel.pem" bin/*.so  ec2-user@ec2-15-228-21-212.sa-east-1.compute.amazonaws.com:/home/ec2-user/bin
scp -r -i "..\..\manoel.pem" k2web/*  ec2-user@ec2-15-228-21-212.sa-east-1.compute.amazonaws.com:/home/ec2-user/k2web
ssh -i "..\..\manoel.pem" ec2-user@ec2-15-228-21-212.sa-east-1.compute.amazonaws.com



