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
	_ "sync"
	"time"
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

//var wg sync.WaitGroup

func main() {

	//runtime.GOMAXPROCS(MaxParallelism())

	var fileMode = flag.Bool("f", false, "If file mode is set to true then it will look for youtube urls serperated by a new line in the files path")
	flag.Var(&urlStrings, "u", "Enter Youtube video url, each url needs the -u command before it")
	flag.Parse()

	if *fileMode == false {
		for _, url := range urlStrings {
			if runtime.GOOS == "windows" {
				//WINDOWS ENVIRONMENT CHECK, TO MAKE SURE THE BINARIES THAT WE ARE USING ARE THE CORRECT ONES
				//TODO
			} else {
				if checkUrl(url) {
					//wg.Add(1)
					//go macDownloader(url)
					macDownloader(url)
				} else {
					fmt.Printf("The url %s is not a proper youtube url\n", url)
				}
			}
		}
		wg.Wait()
		fmt.Printf("Done converting videos\n")
	} else {
		//Load URLS from text file
		//TODO
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

func macDownloader(url string) {
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

	// change the directory to the directory of the videodown
	os.Chdir(youtubeDirectoryPath)
	fmt.Printf("Downloading video %s\n", url)
	cmd := exec.Command("/bin/sh", "-c", "python -m youtube_dl "+url)
	cmd.Run()
	fmt.Printf("Downloading video %s complete\n", url)
	videos := checkExt(".mp4")
	oldVideoPath := filepath.Join(youtubeDirectoryPath, videos[0])
	newVideoPath := filepath.Join(videoDirectoryPath, videos[0])

	// move the file the the vidoes directory
	err = os.Rename(oldVideoPath, newVideoPath)
	if err != nil {
		fmt.Println(err)
	}

	// string magic to ensure the paths are correct, formats, and paths
	// this could probably be cleaned up a bit
	newVideoFileName := strings.Replace(oldVideoPath, videoDirectoryPath, "", -1)
	newVideoFileName = strings.Replace(newVideoFileName, ".mp4", ".mp3", -1)
	oldVideoPath = newVideoPath
	oldVideoPath = strings.Replace(oldVideoPath, oldVideoPath, "\""+oldVideoPath+"\"", -1)
	newVideoPath = filepath.Join(path, mp3Folder)
	newVideoPath = filepath.Join(newVideoPath, newVideoFileName)
	newVideoPath = strings.Replace(newVideoFileName, "youtube-dl-master", "mp3_files", -1)
	newVideoPath = strings.Replace(newVideoPath, newVideoPath, "\""+newVideoPath+"\"", -1)

	//make sure the file in the directory before executing the command
	fmt.Printf("confirming path for %s\n", url)
	os.Chdir(videoDirectoryPath)
	stop := 0
	exit := false
	for {
		videos = checkExt(".mp4")
		if len(videos) > 0 {
			for _, vidName := range videos {
				if strings.Contains(strings.Replace(newVideoFileName, ".mp3", ".mp4", -1), vidName) == true {
					exit = true
					break
				} else {
					fmt.Printf("Video not found\n")
					time.Sleep(1000 * time.Millisecond)
					stop++
				}
			}
		} else {
			fmt.Printf("Nothing in folder\n")
			time.Sleep(1000 * time.Millisecond)
			stop++
		}
		if stop > 15 || exit {
			break
		}
	}
	fmt.Printf("Paths confirmed\n")

	//Ensure path is located where the binaries live
	os.Chdir(path)
	fmt.Printf("Converting video to mp3\n")
	ffmpegCommand := fmt.Sprintf("./ffmpeg -i %s %s", oldVideoPath, newVideoPath)
	out, err := exec.Command("/bin/sh", "-c", ffmpegCommand).CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v %v\n", err, string(out))
	} else {
		fmt.Printf("Removing video\n")
		err = os.Remove(strings.Replace(oldVideoPath, "\"", "", -1))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
	//wg.Done()
}
