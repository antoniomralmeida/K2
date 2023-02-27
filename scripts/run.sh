rm bin/*.bin
version="0.9.0-beta"
build=$(date +%Y%m%d)

git tag  $version 
git push origin --tags

go build  -ldflags "-X 'github.com/antoniomralmeida/k2/pkg/version.Version=$version' -X 'github.com/antoniomralmeida/k2/pkg/version.Build=$build' " -o ./bin/k2web.bin ./cmd/k2web/main.go
go build  -ldflags "-X 'github.com/antoniomralmeida/k2/pkg/version.Version=$version' -X 'github.com/antoniomralmeida/k2/pkg/version.Build=$build' " -o ./bin/k2olivia.bin ./cmd/k2olivia/main.go
go build  -ldflags "-X 'github.com/antoniomralmeida/k2/pkg/version.Version=$version' -X 'github.com/antoniomralmeida/k2/pkg/version.Build=$build' " -o ./bin/k2.bin ./cmd/k2/main.go

rm log/*.json
rm log/*.log
docker ps -aq | xargs docker stop | xargs docker container rm
docker-compose -f ./build/docker-compose.base.yml  build 
docker-compose -f ./build/docker-compose.base.yml  -p "k2" up -d 
./bin/k2olivia.bin &
./bin/k2web.bin &
./bin/k2.bin &
