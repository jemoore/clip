package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type context struct {
	config_file string
	cache_dir   string
	use_cache   bool
	tag         string
	output_file string
	timestamp   string
	duration    string
	videoUrl    string
}

func main() {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	ctx := &context{
		config_file: "input.txt",
		cache_dir:   home_dir + "/Videos/clips/",
		use_cache:   true,
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <tag>")
		fmt.Println("<tag> should match a tag in the input.txt file")
		return
	}
	ctx.tag = os.Args[1]

	ctx.output_file = ctx.cache_dir + ctx.tag + ".mp4"
	if ctx.use_cache {
		// If the output file already exists and if so play it and exit
		_, err := os.Stat(ctx.output_file)
		if err == nil {
			playVideo(ctx.output_file)
			return
		}
	}

	err = getVideoInfo(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = extractVideoClip(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Play the video using a video player package
	playVideo(ctx.output_file)
}

func playVideo(filename string) {
	// Implement a video player functionality or use an external video player package
	cmd := exec.Command("ffplay", filename)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("error playing video: %s : %s", filename, err.Error())
	}
}

func getVideoInfo(ctx *context) error {
	// Open the input file and read info for
	file, err := os.Open(ctx.config_file)
	if err != nil {
		return fmt.Errorf("error opening file: %s : %w", ctx.config_file, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 4 {
			return errors.New("expected 4 parts from config file " +
				"(tag, timestamp, duration, url)")
		}

		if parts[0] == ctx.tag {
			ctx.timestamp = parts[1]
			ctx.duration = parts[2]
			ctx.videoUrl = parts[3]
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %s : %w", ctx.config_file, err)
	}

	if ctx.timestamp == "" || ctx.duration == "" || ctx.videoUrl == "" {
		return fmt.Errorf("error: name %s malformed or not found in input file: %s", ctx.tag, ctx.config_file)
	}
	return nil
}

func extractVideoClip(ctx *context) error {
	// create a temp file to store the downloaded video
	tmp_file, err := createTempFile("clip_video")
	if err != nil {
		return fmt.Errorf("error creating temp file: %s : %s : %w", ctx.tag, tmp_file, err)
	}

	tmp_file = tmp_file + ".mp4"
	defer os.Remove(tmp_file)

	// Download the video using yt-dlp
	cmd := exec.Command("yt-dlp", "-f", "mp4", "-o", tmp_file, ctx.videoUrl)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error downloading video: %s : %s : %w", ctx.tag, ctx.videoUrl, err)
	}

	// Extract the specific portion using ffmpeg
	cmd = exec.Command("ffmpeg", "-i", tmp_file, "-ss",
		ctx.timestamp, "-t", ctx.duration, "-c", "copy", ctx.output_file)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error extracting video clip: %s : %s : %s : %w", ctx.tag, tmp_file, ctx.output_file, err)
	}
	return nil
}

func createTempFile(file_name_pattern string) (string, error) {
	// Create a temporary file in the system's temporary directory
	// although we will not use it directly. Here we test that we can
	// write to the temporary directory and get a base file name.
	tempFile, err := os.CreateTemp(os.TempDir(), file_name_pattern)
	if err != nil {
		return "", err
	}
	// Close the file and return its name
	defer tempFile.Close()
	file_name := tempFile.Name()
	os.Remove(file_name)
	return file_name, nil
}
