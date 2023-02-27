rm bin/*.bin
rm log/*.json

rm log/*.logversion="0.9.0-beta"
build=$(date +%Y%m%d)
git tag  $version 
git push origin --tags
go build  -ldflags "-X 'github.com/antoniomralmeida/k2/pkg/version.Version=$version' -X 'github.com/antoniomralmeida/k2/pkg/version.Build=$build' " -o ./bin/k2web.bin ./cmd/k2web/main.go
go build  -ldflags "-X 'github.com/antoniomralmeida/k2/pkg/version.Version=$version' -X 'github.com/antoniomralmeida/k2/pkg/version.Build=$build' " -o ./bin/k2olivia.bin ./cmd/k2olivia/main.go
go build  -ldflags "-X 'github.com/antoniomralmeida/k2/pkg/version.Version=$version' -X 'github.com/antoniomralmeida/k2/pkg/version.Build=$build' " -o ./bin/k2.bin ./cmd/k2/main.go

docker-compose -f ./build/docker-compose.yml  -p "k2" stop
docker-compose -f ./build/docker-compose.yml  -p "k2" down
docker-compose -f ./build/docker-compose.base.yml  -p "k2" up -d 

gnome-terminal -- ./bin/k2olivia.bin &
gnome-terminal -- ./bin/k2web.bin &
gnome-terminal -- ./bin/k2.bin &
