package version

import "fmt"

var (
	Version string
	Build   string
)

// go build  -ldflags "-X 'github.com/antoniomralmeida/k2/pkg/version.Version=%version%' -X 'github.com/antoniomralmeida/k2/pkg/version.Build=%build%' " -o ./bin/k2web.exe ./cmd/k2web/main.go

func GetVersion() string {
	return Version
}

func GetBuild() string {
	return Build
}

func LogVersion() {
	fmt.Println("version=", Version)
	fmt.Println("build=", Build)

}
