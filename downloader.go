package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
)

var (
	out       bytes.Buffer
	stderr    bytes.Buffer
	usr       *user.User
	videoPath string
	musicPath string
)

// Created so that multiple inputs can be accecpted
type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

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

func checkURL(URL string) bool {
	if strings.Contains(URL, "https://www.youtube.com/watch") || strings.Contains(URL, "https://www.youtube.com/playlist") {
		return true
	}
	return false
}

func downloader(URL string, wg *sync.WaitGroup) {
	log.Println("test", URL)
	defer wg.Done()

}

func main() {
	var urlStrings arrayFlags
	var wg sync.WaitGroup

	var fileMode = flag.Bool("f", false, "If file mode is set to true then it will look for youtube urls serperated by a new line in the files path")
	var formats = flag.String("v", "", "Retrieves video download formats")
	var selectedFormat = flag.String("n", "", "Selected download format")
	var downloadVid = flag.String("d", "", "Download Video")
	flag.Var(&urlStrings, "u", "Enter Youtube video url, each url needs the -u command before it")
	flag.Parse()

	switch {
	case *fileMode:
	case *formats != "":
	case *selectedFormat != "":
	case *downloadVid != "":
	case len(urlStrings) > 0:
		for _, url := range urlStrings {
			if checkURL(url) {
				wg.Add(1)
				go downloader(url, &wg)
			}
		}
		wg.Wait()
	}

}
