package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	_ "strings"
	//"io"
	//"net/http"
	//"net/url"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var urlStrings arrayFlags
var videoFolder = "video_downloads"
var mp3Folder = "mp3_files"
var youtubeFodler = "youtube-dl-master"

func main() {
	// get the absolute path of the script in order to create a video directory if it doesn't exist
	path, err := filepath.Abs("")
	if err != nil {
		fmt.Println("Error locating absulte file paths")
		os.Exit(1)
	}
	videoDirectoryPath := filepath.Join(path, videoFolder)
	exist, err := folderExists(videoDirectoryPath)
	if err != nil {
		fmt.Println("The folder: %s either does not exist or is not in the same directory as make.go", videoDirectoryPath)
		os.Exit(1)
	}
	if !exist {
		os.Mkdir(videoDirectoryPath, 0777)
	}

	mp3DirectoryPath := filepath.Join(path, mp3Folder)
	exist, err = folderExists(mp3DirectoryPath)
	if err != nil {
		fmt.Println("The folder: %s either does not exist or is not in the same directory as make.go", mp3DirectoryPath)
		os.Exit(1)
	}
	if !exist {
		os.Mkdir(mp3DirectoryPath, 0777)
	}

	youtubeDirectoryPath := filepath.Join(path, youtubeFodler)

	var fileMode = flag.Bool("f", false, "If file mode is set to true then it will look for youtube urls serperated by a new line in the files path")
	flag.Var(&urlStrings, "u", "Enter Youtube video url, each url needs the -u command before it")
	flag.Parse()

	// section of code dealing with the downloading of vidoes off of youtube
	// change the directory to the directory of the videodown
	err = os.Chdir(youtubeDirectoryPath)
	if *fileMode == false {
		for _, url := range urlStrings {
			if runtime.GOOS == "windows" {
				fmt.Println("Hello from Windows")
				fmt.Println(url)
			} else {
				// download the video file using the python youtube downloader
				cmd := exec.Command("/bin/sh", "-c", "python -m youtube_dl "+url)
				cmd.Run()
				// move the file the the vidoes directory
				videos := checkExt(".mp4")
				err := os.Rename(filepath.Join(youtubeDirectoryPath, videos[0]), filepath.Join(videoDirectoryPath, videos[0]))
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
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
