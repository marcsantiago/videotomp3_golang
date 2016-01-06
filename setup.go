package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {

	if runtime.GOOS == "windows" {
		path, err := os.Getwd()
		if err != nil {
			fmt.Println("Error locating absolute file paths")
			fmt.Println(err)
			os.Exit(1)
		}
		batScript := filepath.Join(path, "install_ffmpeg.bat")

		fmt.Printf("Copying windows_ffmpeg contents to c:\\FFMPEG and adding the path env c:\\FFMPEG\\bin\n")
		cmd := exec.Command("cmd", "/C", batScript)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

	} else {
		path, err := filepath.Abs("")
		if err != nil {
			fmt.Println("Error locating absolute file paths")
			fmt.Println(err)
			os.Exit(1)
		}
		shellScript := filepath.Join(path, "install_ffmpeg.sh")
		fmt.Printf("Please be patient installing homebrew, ffpmeg, and updating take a while")
		cmd := exec.Command("/bin/sh", shellScript)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

	}
}
