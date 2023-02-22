 go test github.com/antoniomralmeida/k2/internal/models -v -cover -coverprofile test/models.out > test/models.log
 go tool cover -html test/models.out -o test/models.html
 go test -fuzz -v github.com/antoniomralmeida/k2/internal/models 
 