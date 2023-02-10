 go test -v -cover github.com/antoniomralmeida/k2/internal/models -coverprofile cover.out
 go tool cover -html cover.out -o cover.html
