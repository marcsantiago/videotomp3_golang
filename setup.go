package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	path, err := filepath.Abs("")
	if err != nil {
		fmt.Println("Error locating absulte file paths")
		os.Exit(1)
	}

	shellScript := filepath.Join(path, "install_ffmpeg.sh")
	batScript := filepath.Join(path, "install_ffmpeg.bat")

	if runtime.GOOS == "windows" {
		fmt.Printf("Copying windows_ffmpeg contents to c:\\FFMPEG and addinf the path env c:\\FFMPEG\\bin")
		out, err := exec.Command("cmd", "/C", batScript).CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %v %v\n", err, string(out))
		} else {
			fmt.Printf("Homebrew and ffmpeg installed")
		}

	} else {
		fmt.Printf("Please be patient installing homebrew, ffpmeg, and updating take a while")
		out, err := exec.Command("/bin/sh", shellScript).CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %v %v\n", err, string(out))
		} else {
			fmt.Printf("Homebrew and ffmpeg installed")
		}
	}
}
