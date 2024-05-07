package action

import (
	"bytes"
	"fmt"
	"os/exec"
)

// ffmpeg -i rtsp://brun0teixeira:qwerty654321@192.168.30.44:554/stream1 -c:v copy -c:a aac -strict experimental output.mp4

// FIXMEE: this already works so I need to tranform hours to minutes and pass into `-t` flag in ffmpeg
func StartFFMPEGRecording() {
	// Run the ffmpeg command
	// FIXME: this does not work for some reason, the result output.mp4 is corrupted?
	// rand.Seed(time.Now().UnixNano())
	// randomNumber := rand.Intn(101) // generates a random number between 0 and 100
	// a := strconv.Itoa(randomNumber)

	// cmd := exec.CommandContext(ctx, "ffmpeg", "-i", "rtsp://brun0teixeira:qwerty654321@192.168.30.44:554/stream1", "-c:v", "copy", "-c:a", "aac", "-strict", "experimental", "output"+a+".mp4")
	cmd := exec.Command("ffmpeg", "-i", "rtsp://brun0teixeira:qwerty654321@192.168.30.44:554/stream1", "-c:v", "copy", "-c:a", "aac", "-strict", "experimental", "-t", "60", "output.mp4")

	// Capture the output
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	err := cmd.Run()

	// if ctx.Err() == context.DeadlineExceeded {
	// 	fmt.Println("ffmpeg command exceeded the 10-minute timeout")
	// 	// Handle the case where the ffmpeg command exceeds the timeout
	// } else if err != nil {
	// 	fmt.Println("Error running ffmpeg command:", err)
	// 	// Handle other errors, if any
	// } else {
	// 	fmt.Println("ffmpeg command completed within the 10-minute duration")
	// 	// Handle the case where the ffmpeg command completes within the timeout
	// }

	fmt.Printf("ffmpeg output:%s\nffmpeg error:%s\n", output.String(), err.Error())
}
