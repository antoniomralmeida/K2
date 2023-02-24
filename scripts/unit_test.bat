 go test -coverpkg=./... -v -cover -coverprofile test/k2.out > test/k2.log
 go tool cover -html test/k2.out -o test/k2.html
 go test -fuzz -v github.com/antoniomralmeida/k2/internal/models
 