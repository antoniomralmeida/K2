clear
rm bin/*.bin
rm log/*.json
rm log/*.log
rm web/tts/*.mp3


  git config --global user.email "manoel.ribeiro@unilab.edu.br"
  git config --global user.name "Manoel Ribeiro Almeida"
version="0.9.0-beta"
build=$(date +%Y%m%d)
git tag  $version 
git push origin --tags

pkill -f 'k2'
go build  -ldflags "-X 'github.com/antoniomralmeida/k2/pkg/version.Version=$version' -X 'github.com/antoniomralmeida/k2/pkg/version.Build=$build' " -o ./bin/k2web.bin ./cmd/k2web/main.go
go build  -ldflags "-X 'github.com/antoniomralmeida/k2/pkg/version.Version=$version' -X 'github.com/antoniomralmeida/k2/pkg/version.Build=$build' " -o ./bin/k2olivia.bin ./cmd/k2olivia/main.go
go build  -ldflags "-X 'github.com/antoniomralmeida/k2/pkg/version.Version=$version' -X 'github.com/antoniomralmeida/k2/pkg/version.Build=$build' " -o ./bin/k2.bin ./cmd/k2/main.go

#sudo apt install stterm 
terminalpp -e ./bin/k2olivia.bin & 
sleep 10
terminalpp -e ./bin/k2web.bin & 
terminalpp -e ./bin/k2.bin  

pkill -f 'k2'
