package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getVideoFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	exts := map[string]bool{".mp4": true, ".mov": true, ".mkv": true, ".mp3": true, ".m4a": true, ".wav": true}
	var files []string

	if info.IsDir() {
		err := filepath.Walk(path, func(p string, i os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !i.IsDir() && exts[strings.ToLower(filepath.Ext(p))] {
				files = append(files, p)
			}
			return nil
		})
		return files, err
	}

	if exts[strings.ToLower(filepath.Ext(path))] {
		return []string{path}, nil
	}
	return nil, fmt.Errorf("Error: File not supported: use .mp4, .mov, .mkv, .mp3, .m4a or .wav")
}

func requireBinary(name string) string {
	path, err := exec.LookPath(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: '%s' not found in PATH.\n", name)
		fmt.Fprintf(os.Stderr, "Install the prerequisites and try again. See the README for instructions.\n")
		os.Exit(1)
	}
	return path
}

func main() {
	lang := flag.String("lang", "auto", "Language of transcription (e.g., pt, en, es) or 'auto' for automatic detection")
	model := flag.String("model", "", "Path to the Whisper model (.bin). Default: VTTE_MODEL environment variable or ggml-base.bin in the current directory")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: vtte [options] <file_or_directory>")
		fmt.Fprintln(os.Stderr, "\nTranscripts videos .mp4, .mov and .mkv to Markdown using Whisper.")
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  vtte meeting.mp4")
		fmt.Fprintln(os.Stderr, "  vtte -lang pt meeting.mp4")
		fmt.Fprintln(os.Stderr, "  vtte -lang en -model ~/models/ggml-large-v3.bin videos/")
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	modelPath := *model
	if modelPath == "" {
		modelPath = os.Getenv("VTTE_MODEL")
	}
	if modelPath == "" {
		modelPath = "ggml-large-v3.bin"
	}
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: model not found in '%s'.\n", modelPath)
		fmt.Fprintf(os.Stderr, "Download a model em https://huggingface.co/ggerganov/whisper.cpp e use -model to point the path.\n")
		os.Exit(1)
	}

	ffmpegBin := requireBinary("ffmpeg")
	whisperBin := requireBinary("whisper-cli")

	targetPath := flag.Arg(0)
	videos, err := getVideoFiles(targetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading path: %v\n", err)
		os.Exit(1)
	}

	if len(videos) == 0 {
		fmt.Println("No video found.")
		os.Exit(0)
	}

	for _, video := range videos {
			fmt.Printf("\n--- Processing: %s ---\n", filepath.Base(video))

		isWav := strings.ToLower(filepath.Ext(video)) == ".wav"
		var wavFile string
		if isWav {
			wavFile = video
		} else {
			wavFile = video + ".wav"
		}
		txtFile := wavFile + ".txt"

		if !isWav {
			fmt.Println("Extracting audio (ffmpeg)...")
			cmdFfmpeg := exec.Command(ffmpegBin, "-y", "-i", video, "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", wavFile)
			cmdFfmpeg.Stderr = os.Stderr
			if err := cmdFfmpeg.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error extracting audio: %v\n", err)
				continue
			}
		}

		fmt.Printf("Transcripting audio (whisper, language: %s)...\n", *lang)
		cmdWhisper := exec.Command(whisperBin, "-m", modelPath, "-f", wavFile, "-l", *lang, "-otxt")
		cmdWhisper.Stdout = os.Stdout
		cmdWhisper.Stderr = os.Stderr
		if err := cmdWhisper.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error transcribing: %v\n", err)
			os.Remove(wavFile)
			continue
		}

		content, err := os.ReadFile(txtFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading transcription result: %v\n", err)
			if !isWav {
				os.Remove(wavFile)
			}
			continue
		}

		mdContent := fmt.Sprintf("# %s\n\n%s\n", filepath.Base(video), strings.TrimSpace(string(content)))
		mdPath := video[:len(video)-len(filepath.Ext(video))] + ".md"

		if err := os.WriteFile(mdPath, []byte(mdContent), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving markdown: %v\n", err)
		} else {
			fmt.Printf("Saved: %s\n", mdPath)
		}

		if !isWav {
			os.Remove(wavFile)
		}
		os.Remove(txtFile)
	}

	fmt.Println("\nProcessing completed!")
}
