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

	if runtime.GOOS == "windows" {
		//WINDOWS ENVIRONMENT CHECK, TO MAKE SURE THE BINARIES THAT WE ARE USING ARE THE CORRECT ONES
		//TODO
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
