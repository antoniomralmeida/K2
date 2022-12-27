package services

import (
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/antoniomralmeida/k2/initializers"
)

func Run(wg *sync.WaitGroup) {
	//SERVICES WINDOWS
	switch runtime.GOOS {
	case "windows2":
		wd, _ := os.Getwd()
		cmd := wd + "/bin/k2web.exe"
		_, err := exec.Command("cmd.exe", "/c", "start", cmd).Output()
		initializers.Log(err, initializers.Error)
	}
	wg.Done()
}

func Stop() {
	switch runtime.GOOS {
	case "windows":
		cmd := "taskkill /F /IM k2web.exe"
		_, err := exec.Command("cmd.exe", "/c", cmd).Output()
		initializers.Log(err, initializers.Error)
	}
}
