package web

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

func Run(wg *sync.WaitGroup) {
	//WEB SERVER
	fmt.Println(runtime.GOOS)
	switch runtime.GOOS {
	case "windows":
		wd, _ := os.Getwd()
		web := wd + "\\k2web.exe"
		_, err := exec.Command("cmd.exe", "/c", "start", web).Output()
		if err != nil {
			log.Fatal(err)
		}
	case "linux":
		wd, _ := os.Getwd()
		web := wd + "/k2web.bin"
		_, err := exec.Command(web).Output()
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("OS not supported!" + runtime.GOOS)
	}
	wg.Done()
}

func Stop() {
	switch runtime.GOOS {
	case "windows":
		cmd := "taskkill /F /IM k2web.exe"
		_, err := exec.Command("cmd.exe", "/c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("OS not supported!" + runtime.GOOS)
	}
}
