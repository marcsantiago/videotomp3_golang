package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	out       bytes.Buffer
	stderr    bytes.Buffer
	usr       *user.User
	videoPath string
	musicPath string
	path, _   = os.Getwd()
)

// Created so that multiple inputs can be accecpted
type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
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

func checkExt(ext string) []string {
	pathS, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var files []string
	filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	return files
}

func checkURL(URL string) bool {
	if strings.Contains(URL, "https://www.youtube.com/watch") || strings.Contains(URL, "https://www.youtube.com/playlist") {
		return true
	}
	return false
}

func downloader(URL string, wg *sync.WaitGroup, video bool) {
	defer wg.Done()
	defer out.Reset()
	defer stderr.Reset()
	if !video {
		cmd := exec.Command("/usr/local/bin/youtube-dl", "--ignore-errors", "--extract-audio", "--audio-format", "mp3", "-o", "%(title)s.%(ext)s", URL)
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			log.Println("Error Downloading", URL, err)
		} else {
			log.Println(out.String())
			log.Println("")
		}
	} else {
		cmd := exec.Command("/usr/local/bin/youtube-dl", "--ignore-errors", "-f", "bestvideo", URL)
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			log.Println("Error Downloading", URL, err)
		} else {
			log.Println(out.String())
			log.Println("")
		}
	}

	return
}

func moveVids() {
	// move all files after downloading
	videos := checkExt(".m4a")
	videos = append(videos, checkExt(".webm")...)
	videos = append(videos, checkExt(".mp4")...)
	videos = append(videos, checkExt(".3gp")...)
	videos = append(videos, checkExt(".flv")...)
	for _, vid := range videos {
		oldVideoPath := filepath.Join(path, vid)
		newVideoPath := filepath.Join(videoPath, vid)
		os.Rename(oldVideoPath, newVideoPath)
	}
}

func moveMusic() {
	// move all files after downloading
	music := checkExt(".mp3")
	for _, m := range music {
		oldMusicPath := filepath.Join(path, m)
		newMusicPatth := filepath.Join(musicPath, m)
		err := os.Rename(oldMusicPath, newMusicPatth)
		fmt.Println(err, m)
	}
}

func main() {
	var musicStrings arrayFlags
	var videoStrings arrayFlags
	var wg sync.WaitGroup

	var fileMode = flag.Bool("f", false, "If file mode is set to true then it will look for youtube urls serperated by a new line in the files path")
	var fpath = flag.String("p", "", "If file path, needed if fileMode is set to true")
	flag.Var(&videoStrings, "v", "Enter Youtube video url, each url needs the -v command before it")
	flag.Var(&musicStrings, "m", "Enter Youtube video url, each url needs the -m command before it")
	flag.Parse()

	switch {
	case *fileMode:
		if *fpath != "" {
			f, err := os.Open(*fpath)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			var url string
			for scanner.Scan() {
				url = strings.TrimSpace(scanner.Text())
				if checkURL(url) {
					wg.Add(1)
					go downloader(url, &wg, false)
				}
			}
			wg.Wait()
			moveMusic()
		} else {
			fmt.Println("File path not set")
		}
	case len(videoStrings) > 0:
		for _, url := range videoStrings {
			if checkURL(url) {
				wg.Add(1)
				go downloader(url, &wg, true)
			}
		}
		wg.Wait()
		moveVids()
	case len(musicStrings) > 0:
		for _, url := range musicStrings {
			if checkURL(url) {
				wg.Add(1)
				go downloader(url, &wg, false)
			}
		}
		wg.Wait()
		moveMusic()
	}

}
