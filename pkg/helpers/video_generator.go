package helpers

import (
	"fmt"
	"os/exec"
	"strings"
)

func AddAudioToVideo(video_file, audio_file, output_file string) error {
	args := strings.Fields(fmt.Sprint("ffmpeg -y -i ", video_file, " -i ", audio_file, " -map 0:v -map 1:a -c:v copy -shortest ", output_file))
	cmd := exec.Command(args[0], args[1:]...)
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}
