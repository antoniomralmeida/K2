package web

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/antoniomralmeida/k2/initializers"
)

func Run(wg *sync.WaitGroup) {
	//WEB SERVER
	fmt.Println(runtime.GOOS)
	switch runtime.GOOS {
	case "windows":
		wd, _ := os.Getwd()
		web := wd + "/bin/k2web.exe"
		_, err := exec.Command("cmd.exe", "/c", "start", web).Output()
		initializers.Log(err, initializers.Error)
	case "linux":
		wd, _ := os.Getwd()
		web := wd + "/bin/k2web.bin"
		_, err := exec.Command(web).Output()
		initializers.Log(err, initializers.Error)
	default:
		initializers.Log("OS not supported!"+runtime.GOOS, initializers.Error)
	}
	wg.Done()
}

func Stop() {
	switch runtime.GOOS {
	case "windows":
		cmd := "taskkill /F /IM k2web.exe"
		_, err := exec.Command("cmd.exe", "/c", cmd).Output()
		initializers.Log(err, initializers.Error)
	default:
		initializers.Log("OS not supported!"+runtime.GOOS, initializers.Error)
	}
}
