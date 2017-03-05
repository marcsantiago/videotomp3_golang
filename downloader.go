package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	out       bytes.Buffer
	stderr    bytes.Buffer
	usr       *user.User
	videoPath string
	musicPath string
)

func init() {
	defer out.Reset()
	defer stderr.Reset()

	usr, _ = user.Current()

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

	// ensure file system is in place
	parent := filepath.Join(usr.HomeDir, "Desktop/YouTubeFiles")
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		fmt.Printf("The folder %s was created for you\n", parent)
		os.Mkdir(parent, 0755)
	}
	videoPath = filepath.Join(parent, "videos")
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		os.Mkdir(videoPath, 0755)
	}
	musicPath = filepath.Join(parent, "music")
	if _, err := os.Stat(musicPath); os.IsNotExist(err) {
		os.Mkdir(musicPath, 0755)
	}
	return
}

func main() {

}
