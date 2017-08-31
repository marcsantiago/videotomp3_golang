package setup

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func SetupBrew() error {
	cmd := exec.Command("/usr/local/bin/brew", "help")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Homebrew needs to be installed inorder to continue.")
		var input string
		fmt.Println("Would you like to installed homebrew")
		for {
			fmt.Println("Please enter either yes (y) or no (n)")
			fmt.Scanf("%s", &input)
			switch i := input; {
			case strings.EqualFold(i, "y"):
				cmd = exec.Command("/usr/bin/ruby", "-e", "\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\"")
				go io.Copy(cmd.Stdout, os.Stdout)
				go io.Copy(cmd.Stderr, os.Stderr)
				err = cmd.Run()
				if err != nil {
					return errors.New("Error while trying to install homebrew")
				}
				return nil
			case strings.Contains(i, "n"):
				return nil
			default:
				println("Invalid input")
			}
		}
	}
	return nil
}

func SetupYouTubeDL() error {
	var buf bytes.Buffer
	cmd := exec.Command("/usr/local/bin/brew", "list")
	cmd.Stdout = &buf

	err := cmd.Run()
	if err != nil {
		return err
	}

	list := buf.String()
	if !strings.Contains(list, "youtube-dl") {
		var input string
		fmt.Println("youtube-dl is needed, would you like to install it via homebrew")
		for {
			fmt.Println("Please enter either yes (y) or no (n)")
			fmt.Scanf("%s", &input)
			switch i := input; {
			case strings.Contains(i, "y"):
				buf.Reset()
				cmd = exec.Command("/usr/local/bin/brew", "install", "youtube-dl")
				cmd.Stdout = &buf
				err = cmd.Run()
				if err != nil {
					return errors.New("Error while trying to installing youtube-dl")
				}
				return nil
			case strings.Contains(i, "n"):
				return nil
			default:
				fmt.Println("Invalid input")
			}
		}
	}
	return nil
}

func SetupFFMPEG() error {
	var buf bytes.Buffer
	cmd := exec.Command("/usr/local/bin/brew", "list")
	cmd.Stdout = &buf

	err := cmd.Run()
	if err != nil {
		return err
	}

	list := buf.String()
	if !strings.Contains(list, "youtube-dl") {
		var input string
		fmt.Println("youtube-dl is needed, would you like to install it via homebrew")
		for {
			fmt.Println("Please enter either yes (y) or no (n)")
			fmt.Scanf("%s", &input)
			switch i := input; {
			case strings.Contains(i, "y"):
				buf.Reset()
				cmd = exec.Command("/usr/local/bin/brew", "install", "ffmpeg")
				cmd.Stdout = &buf
				err = cmd.Run()
				if err != nil {
					return errors.New("Error while trying to installing ffmpeg")
				}
				return nil
			case strings.Contains(i, "n"):
				return nil
			default:
				fmt.Println("Invalid input")
			}
		}
	}
	return nil
}
