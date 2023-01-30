package lib

import "fmt"

var (
	version string
	build   string
)

func GetVersion() string {
	return version
}

func GetBuild() string {
	return build
}

func LogVersion() {
	fmt.Println("version=", version)
	fmt.Println("build=", build)

}
