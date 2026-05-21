# Video to Text Extract

![Languages used](https://img.shields.io/github/languages/count/isadfrn/vtte?style=flat-square)
![Repository size](https://img.shields.io/github/repo-size/isadfrn/vtte?style=flat-square)
![Last commit](https://img.shields.io/github/last-commit/isadfrn/vtte?style=flat-square)

Cross-platform command-line tool (Windows, Linux, and Mac) that transcribes audio from video and audio files (`.mp4`, `.mov`, `.mkv`, `.mp3`, `.m4a`, `.wav`) into Markdown files, using [Whisper.cpp](https://github.com/ggerganov/whisper.cpp) locally, 100% offline, without sending data to any server.

## Prerequisites

Before using `vtte`, install the following tools on your system:

### 1. FFmpeg

**Mac (Homebrew):**

```bash
brew install ffmpeg
```

**Linux (Debian/Ubuntu):**

```bash
sudo apt install ffmpeg
```

**Windows (Winget):**

```powershell
winget install ffmpeg
```

### 2. Whisper CLI (whisper.cpp)

**Mac (Homebrew):**

```bash
brew install whisper-cpp
```

**Linux / Windows:**
Build from source:

```bash
git clone https://github.com/ggerganov/whisper.cpp
cd whisper.cpp
cmake -B build && cmake --build build --config Release
```

Once built, copy the `whisper-cli` binary (or `whisper-cli.exe`) to a directory that is in your `PATH`.

### 3. Whisper Model

Download the AI model. `ggml-base` is a good starting point (a balance between speed and accuracy):

```bash
# Create a folder for the models
mkdir -p ~/whisper-models

# Download the base model
curl -L https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin \
  -o ~/whisper-models/ggml-base.bin
```

Other available models (from lightest to most accurate):

- `ggml-tiny.bin` — fastest
- `ggml-base.bin` — recommended
- `ggml-small.bin`
- `ggml-medium.bin`
- `ggml-large-v3.bin` — most accurate, requires more memory

## Installing vtte

**Requirement:** [Go 1.22+](https://go.dev/dl/)

```bash
go install github.com/isadfrn/vtte@latest
```

Or build locally:

```bash
git clone https://github.com/isadfrn/vtte
cd vtte
go build -o vtte .
```

## Usage

```
vtte [options] <file_or_directory>
```

### Options

| Flag     | Default         | Description                                                                       |
| -------- | --------------- | --------------------------------------------------------------------------------- |
| `-lang`  | `auto`          | Transcription language (e.g., `pt`, `en`, `es`) or `auto` for automatic detection |
| `-model` | `ggml-base.bin` | Path to the Whisper model file                                                    |

The model can also be set via the `VTTE_MODEL` environment variable.

### Examples

**Transcribe a single file (language detected automatically):**

```bash
vtte meeting.mp4
vtte podcast.mp3
vtte interview.wav
```

**Force Portuguese language:**

```bash
vtte -lang pt meeting.mp4
```

**Use a larger model for higher accuracy:**

```bash
vtte -model ~/whisper-models/ggml-large-v3.bin -lang pt meeting.mp4
```

**Transcribe an entire folder (mixed formats):**

```bash
vtte -lang pt ~/recordings/
```

**Set the model via an environment variable:**

```bash
export VTTE_MODEL=~/whisper-models/ggml-base.bin
vtte folder/with/videos/
```

## Output

For each processed file, `vtte` generates a `.md` file in the same folder as the original:

```
meeting.mp4    →  meeting.md
podcast.mp3    →  podcast.md
interview.wav  →  interview.md
```

The Markdown file contains the title with the video name followed by the transcribed text, ready to be imported into **Google NotebookLM**, **Claude**, or any other AI tool.

## How it works

1. **Audio extraction** — FFmpeg extracts the audio from the video and converts it to PCM 16kHz mono (the ideal format for Whisper)
2. **Transcription** — Whisper CLI processes the audio locally with the chosen model
3. **Markdown** — The text is saved as `.md` with the video name as the title
4. **Cleanup** — Temporary `.wav` and `.txt` files are removed automatically

## Contributing

This repository is using [Gitflow Workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow) and [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/), so if you want to contribute:

- create a branch from develop branch;
- make your contributions;
- open a [Pull Request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request) to develop branch;
- wait for discussion and future approval;

I thank you in advance for any contribution.

## Status

Maintaining

## License

[MIT](./LICENSE)
