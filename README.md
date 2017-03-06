# videotomp3_golang

A mac specific command-line wrapper for [youtube-dl](https://rg3.github.io/youtube-dl/) and [ffmpeg](https://ffmpeg.org/).  This application uses homebrew to install the youtube-dl and ffmpeg binaries to your mac.  If you don't have this software installed it will prompt you to install the software and automate it.  Of course you could download youtube-dl your self and just use that tool alone, but using this wrapper will speed up downloads as it will download videos or mp3s from youtube in parallel.

Oh this will automatically create a folder on your desktop called YouTubeFiles.

**Usage:**  
**Full guide sorry for the details if this is redundant to you

I suggest building a binary.  In order to do so, navigate to the folder with terminal.  Once inside run the command:

    go build downloader.go


This will create the binary that can be used from terminal via the command `./downloader`

**Commands:**
To download one song

    ./downloader -music https://www.youtube.com/watch?v=OcIDeP8_Fto


To download many songs just continue to add urls with the -music command in front

    ./downloader -music https://www.youtube.com/watch?v=OcIDeP8_Fto -music https://www.youtube.com/watch?v=hpFZWeQq_EU

**Same as the music download commands :-)**

To download one video

    ./downloader -video https://www.youtube.com/watch?v=OcIDeP8_Fto


To download many videos just continue to add urls with the -music command in front

    ./downloader -video https://www.youtube.com/watch?v=OcIDeP8_Fto -video https://www.youtube.com/watch?v=hpFZWeQq_EU

Now... that can get tedious so I've also programmed a way to submit urls from a newline delimitated file.  At the moment it only does music downloads.


    ./downloader -path ~/Desktop/myfile.txt -file true


Lastly you can save tons of time by creating a **public** playlist on youtube or finding a public playlist on  youtube and download all the songs or videos from that playlist at once.  Notice the url has the word playlist it in

     ./downloader -music https://www.youtube.com/playlist?list=PL_5Qq5Bm7m2bDQAZjz-Io0zv1Yc__wt2o -playlist true

or... for videos

    ./downloader -video https://www.youtube.com/playlist?list=PL_5Qq5Bm7m2bDQAZjz-Io0zv1Yc__wt2o -playlist true


Feel free to fork, add capabilities, etc, enjoy.
