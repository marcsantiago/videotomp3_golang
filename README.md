**I created this script because...**

I was tired of downloading mp3s from youtube downloader sites. I find them slow (you can only do one link at a time).  By using [youtube-dl](https://rg3.github.io/youtube-dl/) and [ffmpeg](https://www.ffmpeg.org/) I was able to write a script in golang, which takes advantage of the number of cores your pc has to do the downloading in batches. The script checks to see how many cores your computer has in order to download mp3s in parallel.

**Setting Up The Downloader**

The program is design to work on both macs and windows 64 bit.  To get started download the project and extract it's contents or clone the project.  Once you've done that please run the `setup.go` file.  For those of you that do not have golang installed on your computers, I have provided binaries located in their respected folders.  If you choose to use the binaries please copy them to the top level directory (The same directory the setup.go and downloader.go file live).
 If you are a **mac user** the program will attempt to install homebrew, followed by ffmpeg.  If the program fails try running these commands manually:

 1. `ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"`
 2. `brew install ffmpeg --with-fdk-aac --with-ffplay --with-freetype --with-libass --with-libquvi --with-libvorbis --with-libvpx --with-opus --with-x265`
 3. `brew update && brew upgrade ffmpeg`
 
If you are a **windows user** The program will attempt to create a folder in your c drive and copy the contents of the windows_ffmpeg folder to that new directory.  Afterwards it will create a path environment to the ffmpeg bin folder.

**Downloaded Usage**

To create a config file which points to the directories you wish your mp3s and videos to download to please use the command:

    go run downloader.go -c true
You willbe prompted to enter the folder path of both the mp3 folder and video folder you wish to download your contents to.  Example path:

    /User/Desktop/my_mp3s

**Special Note: do not include spaced in your path**
**Correct --> /User/Desktop/my_mp3s Incorrect --> /User/Desktop/my mp3s**


**MP3s**

To download an single mp3 or entire playlist use the command:

    go run downloader.go -u [youtube url]

Example:

    go run downloader.go -u https://www.youtube.com/watch?v=lupDjeumqvU&list=PL7FEE4DBE28ADEB51

    go run downloader.go -u https://www.youtube.com/playlist?list=PL7FEE4DBE28ADEB51

If you wish to download more then one mp3 at once use the command:

    go run downloader.go -u [youtube url] -u [youtube url]...
    
You can also create a text file containing a list of links that you can download from.  Each url should be it's own line.  See the example_download_youtube_urls.txt for a formatting example. Use the command:

    go run downloader.go -f [PATH TO FILE]


**Videos**

Videos can only be downloaded one at a time at the moment. Before you can download a video you need to determine what formats the video can be downloaded in.  To do that use the command:

    go run downloader.go -v [youtube url]
 
Example output:

    [youtube] Setting language
        [youtube] P9pzm5b6FFY: Downloading webpage
        [youtube] P9pzm5b6FFY: Downloading video info webpage
        [youtube] P9pzm5b6FFY: Extracting video information
        [info] Available formats for P9pzm5b6FFY:
        format code extension resolution  note 
        140         m4a       audio only  DASH audio , audio@128k (worst)
        160         mp4       144p        DASH video , video only
        133         mp4       240p        DASH video , video only
        134         mp4       360p        DASH video , video only
        135         mp4       480p        DASH video , video only
        136         mp4       720p        DASH video , video only
        17          3gp       176x144     
        36          3gp       320x240     
        5           flv       400x240     
        43          webm      640x360     
        18          mp4       640x360     
        22          mp4       1280x720    (best)

After you've determined which format you want use this command to download the video:

    go run downloader.go -n 22 -d [youtube url]
    