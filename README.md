# Clip
This application will extract and play a section (clip) from a youtube video.

The video url, start time stamp, and duration are read from a file.
Once the video clip is extracted it is played and stored in a cache location.
If the video clip already exists in the cache location it is immediately played.

The cache location is hard coded to:
```~/Videos/clips```

This application has only been tested on Linux but should work on other platforms
if the external dependencies listed below are available.

## Why
This is just a fun way to capture and view 'meme' video clips.

## Input file
The source video information should be stored in a file named ```input.txt```.
The file should reside at the same location as this application.
The file should contain four entries per line:
```Tag,start timestamp,duration,video url```
The tag should match the name passed to the application and indicates the clip to be played.
```go run main.go rick```

## Possible Improvements
- Allow the cache location to be read from an environment variable and/or the command line.
- Add a command line flag to toggle off the use of the cache location
  - this would allow overwriting the video clip (i.e. a new source url is found)
- Add an option to not permanently store generated clips

## Dependencies
This application calls several external applications to perm tasks:
- yt-dlp - to download the youtube video to clip
- ffmpeg - to extract a portion of the video
- ffplay - to play the video clip