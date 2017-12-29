package downloader

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"sync"

	log "github.com/marcsantiago/logger"
)

const (
	logKey = "Downloader"
)

type Downloader struct {
	wg  sync.WaitGroup
	Buf bytes.Buffer
	Ran bool
}

func (d *Downloader) Add(i int) {
	d.wg.Add(i)
	return
}
func (d *Downloader) Wait() {
	d.wg.Wait()
	d.Ran = true
	return
}

func (d *Downloader) checkURL(URL string) bool {
	if strings.Contains(URL, "https://www.youtube.com/watch") || strings.Contains(URL, "https://www.youtube.com/playlist") {
		return true
	}
	return false
}

func (d *Downloader) Run(URL string, video bool) {
	defer d.wg.Done()
	if !d.checkURL(URL) {
		log.Fatal(logKey, "Url entered was not valid", "url", URL)
	}

	if video {
		cmd := exec.Command("/usr/local/bin/youtube-dl", "--ignore-errors", "-f", "bestvideo", URL)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(logKey, "Video cmd", "error", err)
		}
		return
	}

	cmd := exec.Command("/usr/local/bin/youtube-dl", "--ignore-errors", "--extract-audio", "--audio-format", "mp3", "-o", "%(title)s.%(ext)s", URL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(logKey, "Music cmd", "error", err)
	}
	return
}
