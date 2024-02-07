package imagetovideo

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gocv.io/x/gocv"
)

func makeModule() *VideoManager {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading")
	}
	token := os.Getenv("STABILITY_AI_TOKEN_KEY")
	vm := NewVideoManager("./img", token)
	return vm
}

func TestCanResizeImage(t *testing.T) {
	// given : 유효한 이미지
	// when : 이미지 사이즈 변경 요청
	vm := makeModule()

	resizedFilePath, _ := vm.resizeImage()

	// then : 이미지 저장 파일 확인
	_, err := os.Stat(resizedFilePath)
	assert.False(t, os.IsNotExist(err), "리사이즈된 파일이 정상적으로 생성되지 않았습니다.")
}

// func TestCanPostGenerateVideo(t *testing.T) {
// 	// given : 유효한 토큰 + 이미지
// 	// when : 비디오 생성 요청
// 	vm := makeModule()
// 	id, _ := vm.postGenerateVideo()

// 	// then : 생성 요청 확인(generated_id)
// 	fmt.Println("Generated id : ", id)
// 	assert.True(t, len(id) > 0, "비디오 생성 요청이 정상적으로 수행 되지 않았습니다.")
// }

// func TestCanGetGenerateVideo(t *testing.T) {
// 	// given : generated_id
// 	generatedID := "789de2faf74ef9b52512966bbcc32b23274fec293adea42f1cf995921aacfc3f"

// 	// when : 비디오 전달 요청
// 	vm := makeModule()
// 	generatedVideoPath, _ := vm.getGenerateVideo(generatedID)

// 	// then : 비디오 확인
// 	assert.True(t, filepath.Base(generatedVideoPath) == "generated_video.mp4", "비디오 파일이 정상적으로 생성되지 않았습니다.")
// }

func TestCanMakeReversedVideo(t *testing.T) {
	// given : 유효한 동영상
	// when : 비디오 역재생 영상 요청
	gerneratedVideoPath := "./img/generated_video.mp4"
	reversedVideoPath := "./img/reversed_video.mp4"

	err := os.Remove(reversedVideoPath)
	if err != nil {
		fmt.Println("Fail to delete", err)
	}

	// create video capture input file
	cap, err := gocv.VideoCaptureFile(gerneratedVideoPath)
	if err != nil {
		fmt.Println("Can't capture input file.", err)
	}
	defer cap.Close()

	// get the width, height and FPS of the video frame
	frameWidth := int(cap.Get(gocv.VideoCaptureFrameWidth))
	frameHeight := int(cap.Get(gocv.VideoCaptureFrameHeight))
	fps := int(cap.Get(gocv.VideoCaptureFPS))

	// create video capture out file
	out, err := gocv.VideoWriterFile(reversedVideoPath, "mp4v", float64(fps), frameWidth, frameHeight, true)
	if err != nil {
		fmt.Println("Can't capture output file.", err)
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
	for i := 0; i < len(frames); i++ {
		out.Write(frames[i])
	}

	for i := len(frames) - 1; i >= 0; i-- {
		out.Write(frames[i])
	}
	// then 비디오 확인
}
