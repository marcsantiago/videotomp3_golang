package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	out    bytes.Buffer
	stderr bytes.Buffer
)

func init() {
	defer out.Reset()
	defer stderr.Reset()
	// check that homebrew is installed
	cmd := exec.Command("/usr/local/bin/brew", "help")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Homebrew needs to be installed inorder to continue.")
		var input string
		fmt.Println("Would you like to installed homebrew")
		for {
			fmt.Println("Please enter either yes (y) or no (n)")
			fmt.Scanf("%s", &input)
			switch i := input; {
			case strings.Contains(i, "y"):
				out.Reset()
				stderr.Reset()
				cmd = exec.Command("/usr/bin/ruby", "-e", "\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\"")
				cmd.Stdout = &out
				cmd.Stderr = &stderr
				err = cmd.Run()
				if err != nil {
					panic(errors.New("Error while trying to install homebrew"))
				}
				break
			case strings.Contains(i, "n"):
				os.Exit(0)
			default:
				println("Invalid input")
			}
		}
	}
	// check that youtube-dl and ffmpeg is installed
	out.Reset()
	stderr.Reset()
	cmd = exec.Command("/usr/local/bin/brew", "list")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	l := out.String()
	if !strings.Contains(l, "youtube-dl") {
		var input string
		fmt.Println("youtube-dl is needed, would you like to install it via homebrew")
		for {
			fmt.Println("Please enter either yes (y) or no (n)")
			fmt.Scanf("%s", &input)
			switch i := input; {
			case strings.Contains(i, "y"):
				out.Reset()
				stderr.Reset()
				cmd = exec.Command("/usr/local/bin/brew", "install", "youtube-dl")
				cmd.Stdout = &out
				cmd.Stderr = &stderr
				err = cmd.Run()
				if err != nil {
					panic(errors.New("Error while trying to install youtube-dl"))
				}
				break
			case strings.Contains(i, "n"):
				os.Exit(0)
			default:
				println("Invalid input")
			}
		}
	} else if !strings.Contains(l, "ffmpeg") {
		var input string
		fmt.Println("ffmpeg is needed, would you like to install it via homebrew")
		for {
			fmt.Println("Please enter either yes (y) or no (n)")
			fmt.Scanf("%s", &input)
			switch i := input; {
			case strings.Contains(i, "y"):
				out.Reset()
				stderr.Reset()
				cmd = exec.Command("/usr/local/bin/brew", "install", "ffmpeg")
				cmd.Stdout = &out
				cmd.Stderr = &stderr
				err = cmd.Run()
				if err != nil {
					panic(errors.New("Error while trying to install ffmpeg"))
				}
				break
			case strings.Contains(i, "n"):
				os.Exit(0)
			default:
				println("Invalid input")
			}
		}
	}
	return
}

func main() {

}
