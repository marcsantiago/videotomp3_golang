package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var urlStrings arrayFlags
var youtubeFolder = "youtube-dl-master"
var oldVideoPath string
var newVideoPath string

var mp3DirectoryPath string
var videoDirectoryPath string
var youtubeDirectoryPath string

//var path, _ = filepath.Abs("")
var path string

var wg sync.WaitGroup

func main() {

	if runtime.GOOS == "windows" {
		path, _ = os.Getwd()
	} else {
		path, _ = filepath.Abs("")
	}
	youtubeDirectoryPath = filepath.Join(path, youtubeFolder)

	exist, _ := folderExists("config.json")
	if exist {
		file, err := ioutil.ReadFile("config.json")
		checkFile(err)
		loadedObj := Configs{}
		json.Unmarshal(file, &loadedObj)
		if loadedObj.Mp3Path != "" || loadedObj.Mp3Path != "\n" {
			mp3DirectoryPath = strings.TrimSpace(loadedObj.Mp3Path)
		}
		if loadedObj.VideoPath != "" || loadedObj.VideoPath != "\n" {
			videoDirectoryPath = strings.TrimSpace(loadedObj.VideoPath)
		}
	}

	if mp3DirectoryPath == "" || mp3DirectoryPath == "\n" && exist {
		fmt.Println("Music is being downloaded to the parent path of the downloader.go file inside the mp3_files directory")
		mp3DirectoryPath = filepath.Join(path, "mp3_files")
		exist, _ := folderExists(mp3DirectoryPath)
		if !exist {
			os.Mkdir(mp3DirectoryPath, 0777)
		}
	} else {
		exist, _ := folderExists(mp3DirectoryPath)
		if !exist {
			os.Mkdir(mp3DirectoryPath, 0777)
		}
	}

	if videoDirectoryPath == "" || videoDirectoryPath == "\n" && exist {
		fmt.Println("Music is being downloaded to the parent path of the downloader.go file inside the video_files directory")
		videoDirectoryPath = filepath.Join(path, "video_files")
		exist, _ := folderExists(videoDirectoryPath)
		if !exist {
			os.Mkdir(videoDirectoryPath, 0777)
		}
	} else {
		exist, _ := folderExists(videoDirectoryPath)
		if !exist {
			os.Mkdir(videoDirectoryPath, 0777)
		}
	}

	runtime.GOMAXPROCS(MaxParallelism())

	// SET COMMAND LINE PARSER
	var fileMode = flag.String("f", "false", "If file mode is set to true then it will look for youtube urls serperated by a new line in the files path")
	var Setconfig = flag.Bool("c", false, "Allows the user to create a config.txt file which tells the program where to create mp3 directory")
	var formats = flag.String("v", "", "Retrieves video download formats")
	var selectedFormat = flag.String("n", "", "Selected download format")
	var downloadVid = flag.String("d", "", "Download Video")
	flag.Var(&urlStrings, "u", "Enter Youtube video url, each url needs the -u command before it")

	flag.Parse()

	if *Setconfig == true {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Make sure to please create the folder before entering it's path.")

		fmt.Print("Enter MP3 Folder Path, no spaces in folder name e.g. /User/Desktop/mp3_files: ")
		mp3DirectoryPath, _ := reader.ReadString('\n')
		fmt.Print("Enter Video Folder Path, no spaces in folder name e.g. /User/Desktop/video_files: ")
		videoDirectoryPath, _ := reader.ReadString('\n')

		//NEED TO FIX THIS SECTION
		jsonObj := Configs{
			strings.TrimSpace(mp3DirectoryPath),
			strings.TrimSpace(videoDirectoryPath),
		}

		obj, err := json.Marshal(jsonObj)
		if err != nil {
			fmt.Printf("Error %s\n", err)
		}

		// create a file for writing
		f, err := os.Create("config.json")
		checkFile(err)
		defer f.Close()
		f.Write(obj)
	}

	if *fileMode == "false" {
		if len(urlStrings) > 0 {
			for _, url := range urlStrings {
				if runtime.GOOS == "windows" {
					if checkUrl(url) {
						wg.Add(1)
						go downloadMP3(url, false)
					} else {
						fmt.Printf("The url %s is not a proper youtube url\n The proper prefix is https://www.youtube.com/watch\n", url)
					}
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
			fmt.Printf("Final Steps\n")
			moveMP3s(path, youtubeDirectoryPath)
			fmt.Printf("Check %s for your media\n", mp3DirectoryPath)
		}
	} else {
		f, err := os.Open(*fileMode)
		checkFile(err)
		defer f.Close()

		scanner := bufio.NewScanner(f)
		var url string

		for scanner.Scan() {
			url = strings.TrimSpace(scanner.Text())
			if runtime.GOOS == "windows" {
				if checkUrl(url) {
					wg.Add(1)
					go downloadMP3(url, false)
				} else {
					fmt.Printf("The url %s is not a proper youtube url\n The proper prefix is https://www.youtube.com/watch\n", url)
				}
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
		fmt.Printf("Final Steps\n")
		moveMP3s(path, youtubeDirectoryPath)
		fmt.Printf("Check %s for your media\n", mp3DirectoryPath)
	}

	if *formats != "" {
		url := strings.TrimSpace(*formats)
		if checkUrl(url) {
			if runtime.GOOS == "windows" {
				checkVideoFormats(url, false)
			} else {
				checkVideoFormats(url, true)
			}
		} else {
			fmt.Printf("The url %s is not a proper youtube url\n The proper prefix is https://www.youtube.com/watch\n", url)
		}
	}

	if *downloadVid != "" && *selectedFormat != "" {
		n := strings.TrimSpace(*selectedFormat)
		u := strings.TrimSpace(*downloadVid)
		if _, err := strconv.Atoi(n); err == nil {
			if checkUrl(u) {
				url := fmt.Sprintf("%s %s", n, u)
				if runtime.GOOS == "windows" {
					downloadVideo(url, false)
				} else {
					downloadVideo(url, true)
				}
			} else {
				fmt.Printf("Do did enter a properly formatted link.  Example Command: go run downloader.go -d 22 https://www.youtube.com/watch\n")
			}
		} else {
			fmt.Printf("Do did enter a properly formatted link.  Example Command: go run downloader.go -d 22 https://www.youtube.com/watch\n")
		}
	}
	fmt.Printf("Done downloadning video\n")
	fmt.Printf("Final Steps\n")
	moveVids(path, youtubeDirectoryPath)
	fmt.Printf("Check %s for your media\n", videoDirectoryPath)
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

type Configs struct {
	Mp3Path   string
	VideoPath string
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
			if err == nil && r && (strings.Contains(f.Name(), "pyc") == false || strings.Contains(f.Name(), "mp3.py") == false) {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	return files
}

func checkUrl(url string) bool {
	if strings.Contains(url, "https://www.youtube.com/watch") == true || strings.Contains(url, "https://www.youtube.com/playlist") == true {
		return true
	}
	return false
}

func checkFile(e error) {
	if e != nil {
		panic(e)
	}
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
		// change the directory to the directory of the youtube-dl
		os.Chdir(youtubeDirectoryPath)
		fmt.Printf("Downloading mp3 %s\n", url)

		cmd := exec.Command("/bin/sh", "-c", "python -m  youtube_dl --ignore-errors --extract-audio --audio-format mp3 -o \"%(title)s.%(ext)s \" "+url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

	} else {
		// change path to top level where youtube-dl.exe lives
		os.Chdir(path)
		tool := fmt.Sprintf("youtube-dl.exe --ignore-errors --extract-audio --audio-format mp3 -o \"%%(title)s.%%(ext)s \" " + url)
		cmd := exec.Command("cmd", "/C", tool)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
	wg.Done()
}

func checkVideoFormats(url string, mac bool) {
	if mac == true {
		// change the directory to the directory of the youtube-dl
		os.Chdir(youtubeDirectoryPath)
		fmt.Printf("Downloading video %s\n", url)

		cmd := exec.Command("/bin/sh", "-c", "python -m  youtube_dl -F "+url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

	} else {
		// change path to top level where youtube-dl.exe lives
		os.Chdir(path)
		cmd := exec.Command("cmd", "/C", "python -m  youtube_dl -F "+url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

func downloadVideo(url string, mac bool) {
	if mac == true {
		// change the directory to the directory of the youtube-dl
		os.Chdir(youtubeDirectoryPath)
		fmt.Printf("Downloading video %s\n", url)

		cmd := exec.Command("/bin/sh", "-c", "python -m  youtube_dl -f "+url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

	} else {
		// change path to top level where youtube-dl.exe lives
		os.Chdir(path)
		cmd := exec.Command("cmd", "/C", "python -m  youtube_dl -f "+url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

func moveMP3s(path string, youtubeDirectoryPath string) {
	// move all files after downloading
	if runtime.GOOS == "windows" {
		os.Chdir(path)
	} else {
		os.Chdir(youtubeDirectoryPath)
	}

	videos := checkExt(".mp3")

	for _, vid := range videos {
		if runtime.GOOS == "windows" {
			oldVideoPath = filepath.Join(path, vid)
			newVideoPath = filepath.Join(mp3DirectoryPath, vid)
			newVideoPath = strings.Replace(newVideoPath, "#", "", -1)

		} else {
			oldVideoPath = filepath.Join(youtubeDirectoryPath, vid)
			newVideoPath = filepath.Join(mp3DirectoryPath, vid)
		}
		// move the file the the vidoes directory
		os.Rename(oldVideoPath, newVideoPath)
	}
}

func moveVids(path string, youtubeDirectoryPath string) {
	// move all files after downloading
	if runtime.GOOS == "windows" {
		os.Chdir(path)
	} else {
		os.Chdir(youtubeDirectoryPath)
	}

	videos := checkExt(".m4a")
	videos = append(videos, checkExt(".webm")...)
	videos = append(videos, checkExt(".mp4")...)
	videos = append(videos, checkExt(".3gp")...)
	videos = append(videos, checkExt(".flv")...)

	for _, vid := range videos {
		if runtime.GOOS == "windows" {
			oldVideoPath = filepath.Join(path, vid)
			newVideoPath = filepath.Join(videoDirectoryPath, vid)
			//newVideoPath = strings.Replace(newVideoPath, "#", "", -1)

		} else {
			oldVideoPath = filepath.Join(youtubeDirectoryPath, vid)
			newVideoPath = filepath.Join(videoDirectoryPath, vid)
		}
		// move the file the the vidoes directory
		os.Rename(oldVideoPath, newVideoPath)
	}
}
