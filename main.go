package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/marcsantiago/logger"
	"github.com/marcsantiago/videotomp3_golang/internal/downloader"
	pl "github.com/marcsantiago/videotomp3_golang/internal/playlist"
	"github.com/marcsantiago/videotomp3_golang/internal/setup"
)

var (
	usr       *user.User
	path, _   = os.Getwd()
	videoPath string
	musicPath string
)

const (
	logKey = "Main"
)

// Created so that multiple inputs can be accecpted
type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}

func init() {
	err := setup.SetupBrew()
	if err != nil {
		log.Fatal(logKey, "Issue with brew setup", "error", err)
	}

	err = setup.SetupYouTubeDL()
	if err != nil {
		log.Fatal(logKey, "Issue with youtubedl setup", "error", err)
	}

	err = setup.SetupFFMPEG()
	if err != nil {
		log.Fatal(logKey, "Issue with ffmpeg setup", "error", err)
	}

	usr, _ = user.Current()

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

func checkExt(ext string) (files []string) {
	pathS, err := os.Getwd()
	if err != nil {
		log.Fatal(logKey, "video cmd", "error", err)
	}
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
		os.Rename(oldMusicPath, newMusicPatth)
	}
}

func main() {
	var musicStrings, videoStrings arrayFlags
	var downloader downloader.Downloader

	var fileMode = flag.Bool("file", false, "If file mode is set to true then it will look for youtube urls serperated by a new line in the files path")
	var fpath = flag.String("path", "", "If file path, needed if fileMode is set to true")
	var playlist = flag.Bool("playlist", false, "Download a playlist faster")
	flag.Var(&videoStrings, "video", "Enter Youtube video url, each url needs the -video command before it")
	flag.Var(&musicStrings, "music", "Enter Youtube music url, each url needs the -music command before it")
	flag.Parse()

	conf := struct {
		fileMode     bool
		filePath     string
		playlistMode bool
		videos       arrayFlags
		music        arrayFlags
	}{
		fileMode:     *fileMode,
		filePath:     *fpath,
		playlistMode: *playlist,
		videos:       videoStrings,
		music:        musicStrings,
	}

	list := struct {
		urls  []string
		video bool
	}{}
	if conf.fileMode {
		if conf.filePath == "" {
			log.Fatal(logKey, "File path not set")
		}
		f, err := os.Open(*fpath)
		if err != nil {
			log.Fatal(logKey, "Issue opening file", "error", err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			url := strings.TrimSpace(scanner.Text())
			list.urls = append(list.urls, url)
		}
		if scanner.Err() != nil {
			log.Fatal(logKey, "Scanning", "error", err)
		}
	}

	if conf.playlistMode {
		if len(conf.videos) < 1 && len(conf.music) < 1 {
			fmt.Println("Please enter a single video url -video myurl  or music url -music myurl in conjuction with this command")
			return
		}
		var item pl.PlayList
		var u string
		if len(conf.videos) == 1 {
			u = conf.videos[0]
			list.video = true
		} else if len(conf.music) == 1 {
			u = conf.music[0]
		}

		cmd := exec.Command("/usr/local/bin/youtube-dl", u, "--flat-playlist", "--dump-single-json")
		cmd.Stdout = &downloader.Buf
		err := cmd.Run()
		if err != nil {
			log.Fatal(logKey, "Issue getting playlist json", "error", err)
		}

		err = json.Unmarshal(downloader.Buf.Bytes(), &item)
		if err != nil {
			log.Fatal(logKey, "Issue with json unmarshal", "error", err)
		}

		for _, entry := range item.Entries {
			list.urls = append(list.urls, fmt.Sprintf("https://www.youtube.com/watch?v=%s", entry.URL))
		}
	}

	if len(conf.videos) > 0 {
		list.video = true
		for _, url := range conf.videos {
			list.urls = append(list.urls, url)
		}
	}

	if len(conf.music) > 0 {
		for _, url := range conf.music {
			list.urls = append(list.urls, url)
		}
	}

	if len(list.urls) > 0 {
		for _, url := range list.urls {
			downloader.Add(1)
			go downloader.Run(url, list.video)
		}
		downloader.Wait()
		if list.video {
			moveVids()
			return
		}
		moveMusic()
	}

	if !downloader.Ran {
		flag.PrintDefaults()
	}
	return
}
