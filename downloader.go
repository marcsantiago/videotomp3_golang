package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
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

type VideoFilePaths struct {
	files []string
}

func (v *VideoFilePaths) add(filepath string) []string {
	v.files = append(v.files, filepath)
	return v.files
}

var urlStrings arrayFlags
var videoFolder = "video_downloads"
var mp3Folder = "mp3_files"
var youtubeFolder = "youtube-dl-master"
var ffmpeg = "ffmpeg-2.8.4"

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

	youtubeDirectoryPath := filepath.Join(path, youtubeFolder)

	var fileMode = flag.Bool("f", false, "If file mode is set to true then it will look for youtube urls serperated by a new line in the files path")
	flag.Var(&urlStrings, "u", "Enter Youtube video url, each url needs the -u command before it")
	flag.Parse()

	// section of code dealing with the downloading of videos off of youtube
	// create a list of file names of vidoes
	//var videofiles []string
	//var videofiles VideoFilePaths
	if *fileMode == false {
		for _, url := range urlStrings {
			// change the directory to the directory of the videodown
			os.Chdir(youtubeDirectoryPath)
			if runtime.GOOS == "windows" {
				//TODO
			} else {
				// download the video file using the python youtube downloader
				cmd := exec.Command("/bin/sh", "-c", "python -m youtube_dl "+url)
				cmd.Run()

				videos := checkExt(".mp4")
				oldVideoPath := filepath.Join(youtubeDirectoryPath, videos[0])
				newVideoPath := filepath.Join(videoDirectoryPath, videos[0])

				// move the file the the vidoes directory
				err := os.Rename(oldVideoPath, newVideoPath)
				//videofiles.add(filepath.Join(videoDirectoryPath, videos[0]))
				if err != nil {
					fmt.Println(err)
				}

				// change the working directory to were ffmpeg lives
				os.Chdir(filepath.Join(path, ffmpeg))
				//remove the path from the movie name and change it's path
				newVideoFileName := strings.replace(oldVideoPath, videoDirectoryPath, "", -1)
				newVideoFileName = strings.replace(newVideoFileName, ".mp4", ".mp3", -1)
				oldVideoPath = newVideoPath
				newVideoPath = filepath.Join(path, mp3Folder)
				newVideoPath = filepath.Join(newVideoPath, newVideoFileName)

				fmt.Println(oldVideoPath)
				fmt.Println(newVideoPath)
				cmd = exec.Command("/bin/sh", "-c", "./ffmpeg -i %s %s", oldVideoPath, newVideoPath)
				err = cmd.Run()
				fmt.Println(err)
			}
		}
	} else {
		//TODO
	}
	//test
	// videofiles.add("/Users/marcsantiago/Desktop/videotomp3_golang/video_downloads/Day\\ 24\\ -\\ Kendall\\ Jenner\\ by\\ James\\ Lima\\ \\ \\(LOVE\\ Advent\\ 2015\\)-AmeSgBd-KVE.mp4 ")
	//   exit status 127
	//   //ffmpeg -i filename.mp4 filename.mp3
	//   for _, videoFile := range videofiles.files {
	//     vfile := mp3Folder + strings.Replace(videoFile, path, "", -1)
	//     fmt.Println(vfile)
	//     cmd := exec.Command("./ffmpeg -i", videoFile+" "+vfile)
	//     err := cmd.Run()
	//     fmt.Println(err)
	//   }

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

func confirmUrl(url string) (bool, string) {

}
