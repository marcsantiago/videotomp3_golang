package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	_ "time"
)

var urlStrings arrayFlags
var mp3Folder = "mp3_files"
var youtubeFolder = "youtube-dl-master"

var wg sync.WaitGroup

func main() {

	runtime.GOMAXPROCS(MaxParallelism())

	var fileMode = flag.String("f", "false", "If file mode is set to true then it will look for youtube urls serperated by a new line in the files path")
	flag.Var(&urlStrings, "u", "Enter Youtube video url, each url needs the -u command before it")
	flag.Parse()

	if *fileMode == "false" {
		for _, url := range urlStrings {
			if runtime.GOOS == "windows" {
				//WINDOWS ENVIRONMENT CHECK, TO MAKE SURE THE BINARIES THAT WE ARE USING ARE THE CORRECT ONES
				//TODO
			} else {
				if checkUrl(url) {
					wg.Add(1)
					go downloadMP3(url, true)
				} else {
					fmt.Printf("The url %s is not a proper youtube url\n The proper prefix is https://www.youtube.com/watch\n", url)
				}
			}
		}
		wg.Wait()
		fmt.Printf("Done converting videos\n")
	} else {
		if runtime.GOOS == "windows" {
			//WINDOWS ENVIRONMENT CHECK, TO MAKE SURE THE BINARIES THAT WE ARE USING ARE THE CORRECT ONES
			//TODO
		} else {
			f, err := os.Open(*fileMode)
			if err != nil {
				fmt.Println(err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			var url string
			for scanner.Scan() {
				url = scanner.Text()
				if checkUrl(url) {
					wg.Add(1)
					go downloadMP3(url, true)
				} else {
					fmt.Printf("The url %s is not a proper youtube url\n The proper prefix is https://www.youtube.com/watch\n", url)
				}
			}
			wg.Wait()
		}
	}

}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type VideoFilePaths struct {
	files []string
}

func (v *VideoFilePaths) add(filepath string) []string {
	v.files = append(v.files, filepath)
	return v.files
}

func folderExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
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

func checkUrl(url string) bool {
	if strings.Contains(url, "https://www.youtube.com/watch") == true {
		return true
	}
	return false
}

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func downloadMP3(url string, mac bool) {
	if mac == true {
		path, err := filepath.Abs("")
		if err != nil {
			fmt.Println("Error locating absulte file paths")
			os.Exit(1)
		}
		youtubeDirectoryPath := filepath.Join(path, youtubeFolder)
		mp3DirectoryPath := filepath.Join(path, mp3Folder)
		//create mp3 dicrectory
		exist, err := folderExists(mp3DirectoryPath)
		if err != nil {
			fmt.Println("The folder: %s either does not exist or is not in the same directory as make.go", mp3DirectoryPath)
			os.Exit(1)
		}
		if !exist {
			os.Mkdir(mp3DirectoryPath, 0777)
		}

		// change the directory to the directory of the youtube-dl
		os.Chdir(youtubeDirectoryPath)

		fmt.Printf("Downloading video %s\n", url)

		out, err := exec.Command("/bin/sh", "-c", "python -m  youtube_dl --extract-audio --audio-format mp3 -o \"%(title)s.%(ext)s \" "+url).CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %v %v\n", err, string(out))
		} else {
			fmt.Printf("Downloading video %s complete\n", url)
		}
		// move the videos from the youtube-dl folder to the mp3_files folder
		videos := checkExt(".mp3")
		var oldVideoPath string
		var newVideoPath string
		for _, vid := range videos {
			oldVideoPath = filepath.Join(youtubeDirectoryPath, vid)
			newVideoPath = filepath.Join(mp3DirectoryPath, vid)
			// move the file the the vidoes directory
			os.Rename(oldVideoPath, newVideoPath)

		}

	} else {
		//WINDOWS
		//TODO
	}
	wg.Done()
}
