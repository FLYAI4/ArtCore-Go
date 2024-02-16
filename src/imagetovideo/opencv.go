package imagetovideo

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gocv.io/x/gocv"
)

// Make long video with openCV
func (vm *VideoManager) makeReversedVideo() (string, error) {
	inputFilePath := filepath.Join(vm.userFolderPath, "generated_video.mp4")
	outFilePath := filepath.Join(vm.userFolderPath, "reversed_video.mp4")
	deleteFile(outFilePath)

	// wait inputFilePath
	for {
		_, err := os.Stat(inputFilePath)
		if err == nil {
			fmt.Printf("Find file!!: %s\n", inputFilePath)
			break
		}
		fmt.Println("The file is not created. try again after 3 sec.")
		time.Sleep(3 * time.Second)
	}

	// create video capture input file
	cap, err := gocv.VideoCaptureFile(inputFilePath)
	if err != nil {
		fmt.Println("Can't capture input file.", err)
		return "", err
	}
	defer cap.Close()

	// get the width, height and FPS of the video frame
	frameWidth := int(cap.Get(gocv.VideoCaptureFrameWidth))
	frameHeight := int(cap.Get(gocv.VideoCaptureFrameHeight))
	fps := int(cap.Get(gocv.VideoCaptureFPS))

	// create video capture out file
	out, err := gocv.VideoWriterFile(outFilePath, "mp4v", float64(fps), frameWidth, frameHeight, true)
	if err != nil {
		fmt.Println("Can't capture output file.", err)
		return "", err
	}
	defer out.Close()

	var frames []gocv.Mat
	for {
		frame := gocv.NewMat()
		if ok := cap.Read(&frame); !ok {
			break
		}
		defer frame.Close()

		frames = append(frames, frame.Clone())
	}

	// write forward play
	for i := 0; i < len(frames); i++ {
		out.Write(frames[i])
		out.Write(frames[i])
	}

	// write backward play
	for i := len(frames) - 1; i >= 0; i-- {
		out.Write(frames[i])
		out.Write(frames[i])
	}

	return outFilePath, nil
}
