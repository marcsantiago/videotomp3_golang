@echo on

MKDIR c:\FFMPEG

XCOPY windows_ffmpeg c:\FFMPEG /y /s
SETX PATH "%PATH%;c:\FFMPEG\bin" 