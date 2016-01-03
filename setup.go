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
		out, err := exec.Command("cmd", "/C", batScript).CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %v %v\n", err, string(out))
		} else {
			fmt.Printf("Done installing ffmpeg")
		}

	} else {
		path, err := filepath.Abs("")
		if err != nil {
			fmt.Println("Error locating absolute file paths")
			fmt.Println(err)
			os.Exit(1)
		}
		shellScript := filepath.Join(path, "install_ffmpeg.sh")
		fmt.Printf("Please be patient installing homebrew, ffpmeg, and updating take a while")
		out, err := exec.Command("/bin/sh", shellScript).CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %v %v\n", err, string(out))
		} else {
			fmt.Printf("Homebrew and ffmpeg installed")
		}
	}
}
