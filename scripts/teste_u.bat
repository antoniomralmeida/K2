 go test github.com/antoniomralmeida/k2/internal/models -cover -coverprofile test/cover.out
 go tool cover -html test/cover.out -o test/cover.html
 cd internal\models
 rem go test -fuzz .  -v -timeout 10s
cd ..\..
