Mac Users:
-Install homebrew and ffmpeg on your mac first

Paste this into your terminal:

ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)" && brew install ffmpeg --with-fdk-aac --with-ffplay --with-freetype --with-libass --with-libquvi --with-libvorbis --with-libvpx --with-opus --with-x265 && brew update && brew upgrade ffmpeg

OR

Run the make.go file