 go test github.com/antoniomralmeida/k2/internal/models -cover -coverprofile cover.out
 go tool cover -html cover.out -o cover.html
 cd internal\models
 go test -fuzz .  -v -timeout 10s
cd ..\..
