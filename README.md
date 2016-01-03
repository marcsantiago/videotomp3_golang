<p><b>Reason for creation:</b></p>
<p>
  I was tired of downloading mp3s from youtube on sites that allow you to input one link at a time.  
  By using youtube-dl and ffmpeg I was able to get around that.  
  The go program will download many mp3s in parallel.  The script checks to see how many cores you have and using a gorountine, downloads the mp3s.
</p>

<p><b>Installation:</b></p>
<p>
  The program is design to work on both macs and windows 64 bit.  
  The get started download the project and extract it's contents.  
  From your command line, type <code>go run setup.go</code>.  
  If you are a <b>mac user</b> the program will attempt to install homebrew, followed by the installtion of ffmpeg.  
  If the program files try to manually installed using these three commands
  <ol>
    <li><code>ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"</code></li>
    <li><code>brew install ffmpeg --with-fdk-aac --with-ffplay --with-freetype --with-libass --with-libquvi --with-libvorbis --with-libvpx --with-opus --with-x265</code></li>
    <li><code>brew update && brew upgrade ffmpeg</code></li>
  </ol>
  If you are a <b>windows user</b> The program will attempt to create a folder in your c drive and copy the contents of the windows_ffmpeg folder to that new directory.  Afterwards it will create a path environment to the it's bin folder.
</p>
<p><b>Notes:</b></p>
<p>
  If you wish can can build the downloader using <code>go run downloader.go</code>, which will create a binary of the script
</p>
<p>
  If you are a windows user the program depends on the on the youtube-dl.exe being in the same directory as the downloader.go file or binary.
</p>
<p>
  If you are a mac user the program depends on the on the youtube-dl-master being in the same directory as the downloader.go file or binary.
</p>
<p>
  I suggest you create a shortcut or alias to the binary or go file.
</p>
<p><b>Setup Directory:</b></p>
<p>
  As a default mp3s are downloaded to the mp3_files directory located in the same directory as the downloader script.  
  To change the setting run:
  <code>go run downloader.go -c true [PATH TO DIRECTORY YOU WISH MUSIC TO SAVE TO] </code>
  That will create a config.txt file which contains the path to the directory, you can change the path by running the command again.
  The go program will attempt to create the directory if it doesn't already exist
</p>

<p><b>Usage</b></p>
<p>
  To download mp3s use the command:
  <code>go run downloader.go -u [youtube url]</code>
  If you wish to downloader more then one mp3:
  <code>go run downloader.go -u [youtube url] -u [youtube url]..</code>
</p>
<p>
  You can also specify the path of a text document containing youtube urls that are newline deliminated:
  <code>go run downloader.go -f [PATH TO FILE]</code>
</p>
<p><b>Note</b></p>
<p>
  The program expects the youtube urls to be formatted as such: https://www.youtube.com/watch?v=Ek0SgwWmF9w
</p>