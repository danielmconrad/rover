package rover

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

// Frame NEEDSCOMMENT
type Frame []byte

var (
	raspividArgs = "raspivid -t 0 -w %d -h %d -fps %d -pf baseline -o -"

	ffmpegArgsMac = []string{
		"ffmpeg",
		"-f", "avfoundation", "-framerate", "30", "-pixel_format", "yuyv422",
		"-video_size", "640x480", "-i", "0",
		"-vcodec", "libx264", "-profile:v", "baseline", "-pix_fmt", "yuv420p", "-level:v", "4.2",
		"-preset", "ultrafast", "-tune", "zerolatency", "-bufsize", "0", "-crf", "22",
		"-f", "rawvideo", "-",
	}

	initialFrameCount = 4
	nalSeparator      = []byte{0x00, 0x00, 0x00, 0x01}
	isMac             = runtime.GOOS == "darwin"
)

// StartCamera NEEDSCOMMENT
func StartCamera(ctx context.Context, width, height, framerate uint64) (chan Frame, []Frame) {
	frameChan := make(chan Frame)
	initialFrames := []Frame{}

	args := strings.Split(fmt.Sprintf(raspividArgs, width, height, framerate), " ")

	if isMac {
		args = ffmpegArgsMac
	}

	logInfo("Starting video: ", args)
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logPanic("stdout error", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logPanic("stderr error", err)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(Frame{}, 1024*1024)
	scanner.Split(splitAtNALSeparator)

	go func() {
		logInfo("Scanning video")
		for scanner.Scan() {
			select {
			case frameChan <- scanner.Bytes():
			case <-ctx.Done():
				return
			}
		}
	}()

	var initialFramesWait sync.WaitGroup
	initialFramesWait.Add(1)
	go func() {
		defer initialFramesWait.Done()

		logInfo("Collecting initial frames")

		for frame := range frameChan {
			if len(initialFrames) >= initialFrameCount {
				logInfo("All initial frames collected")
				return
			}

			initialFrames = append(initialFrames, frame)
			logInfo("Frame collected")
		}
	}()

	go func() {
		if err := cmd.Start(); err != nil {
			logPanic("command start error", err)
		}

		if err := cmd.Wait(); err != nil {
			stderrLog, _ := ioutil.ReadAll(stderr)
			logPanic("command wait error", err, stderrLog)
		}
	}()

	initialFramesWait.Wait()

	return frameChan, initialFrames
}

func splitAtNALSeparator(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if idx := bytes.Index(data, nalSeparator); idx >= 0 {
		if idx == 0 {
			return idx + 4, nil, nil
		}
		return idx + 1, append(nalSeparator, data[:idx]...), nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}
