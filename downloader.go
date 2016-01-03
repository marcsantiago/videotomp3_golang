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
)

var urlStrings arrayFlags
var youtubeFolder = "youtube-dl-master"
var oldVideoPath string
var newVideoPath string

var mp3DirectoryPath string
var youtubeDirectoryPath string

//var path, _ = filepath.Abs("")
var path string


var wg sync.WaitGroup

func main() {

	if runtime.GOOS == "windows"{
		path, _ = os.Getwd()
	} else {
		path, _ = filepath.Abs("")
	}

	exist, _ := folderExists("config.txt")
	if exist {
		f, err := os.Open("config.txt")
		checkFile(err)
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" || line == "\n" {
				continue
			} else {
				mp3DirectoryPath = strings.TrimSpace(line)
			}
		}
	}

	

	youtubeDirectoryPath = filepath.Join(path, youtubeFolder)

	runtime.GOMAXPROCS(MaxParallelism())

	// SET COMMAND LINE PARSER
	var fileMode = flag.String("f", "false", "If file mode is set to true then it will look for youtube urls serperated by a new line in the files path")
	var Setconfig = flag.Bool("c", false, "Allows the user to create a config.txt file which tells the program where to create mp3 directory")

	flag.Var(&urlStrings, "u", "Enter Youtube video url, each url needs the -u command before it")
	flag.Parse()

	if *Setconfig == true {
		// create a file for writing
		f, err := os.Create("config.txt")
		checkFile(err)
		defer f.Close()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Folder Path, e.g. /User/Desktop/mp3_files: ")
		mp3DirectoryPath, _ := reader.ReadString('\n')
		exist, err := folderExists(mp3DirectoryPath)
		if err != nil {
			fmt.Printf("Error: %s \n", err)
		}
		if !exist {
			fmt.Printf("The folder: %s either does not exist creating it now\n", mp3DirectoryPath)
			os.Mkdir(mp3DirectoryPath, 0777)
			file, err := f.WriteString(mp3DirectoryPath + "\n")
			checkFile(err)
			fmt.Printf("wrote %d bytes\n", file)
			f.Sync()
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
			fmt.Printf("Check %s for your media\n", mp3DirectoryPath)
			fmt.Printf("Done converting videos\n")
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
					fmt.Printf("Check %s for your media\n", mp3DirectoryPath)
					fmt.Printf("The url %s is not a proper youtube url\n The proper prefix is https://www.youtube.com/watch\n", url)
				}
			}
		}
		wg.Wait()
	}

	// move all files after downloading
	if runtime.GOOS == "windows" {
		os.Chdir(path)
	} else {
		os.Chdir(youtubeDirectoryPath)
	}

	videos := checkExt(".mp3")
	fmt.Println("before for loop")
	for _, vid := range videos {
		
		if runtime.GOOS == "windows" {
			oldVideoPath = filepath.Join(path, vid)
			newVideoPath = filepath.Join(mp3DirectoryPath, vid)
		} else{
			oldVideoPath = filepath.Join(youtubeDirectoryPath, vid)
			newVideoPath = filepath.Join(mp3DirectoryPath, vid)	
		}
		
		// move the file the the vidoes directory
		os.Rename(oldVideoPath, newVideoPath)
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
			if err == nil && r && (strings.Contains(f.Name(), "pyc") == false || strings.Contains(f.Name(), "mp3.py") == false){
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
		fmt.Printf("Downloading video %s\n", url)

		cmd := exec.Command("/bin/sh", "-c", "python -m  youtube_dl --extract-audio --audio-format mp3 -o \"%(title)s.%(ext)s \" "+url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

	} else {
		// change path to top level where youtube-dl.exe lives
		os.Chdir(path)
		tool := fmt.Sprintf("youtube-dl.exe --extract-audio --audio-format mp3 "+ url)
		cmd := exec.Command("cmd", "/C", tool)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
	wg.Done()
}
