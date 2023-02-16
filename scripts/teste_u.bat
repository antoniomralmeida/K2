 go test github.com/antoniomralmeida/k2/internal/models -v -cover -coverprofile test/models.out > test/models.log
 go tool cover -html test/models.out -o test/models.html
 cd internal\models
 rem go test -fuzz .  -v -timeout 10s
cd ..\..
