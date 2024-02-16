package imagetovideo

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
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

func TestCanPostGenerateVideo(t *testing.T) {
	// given : 유효한 토큰 + 이미지
	// when : 비디오 생성 요청
	vm := makeModule()
	id, _ := vm.postGenerateVideo()

	// then : 생성 요청 확인(generated_id)
	fmt.Println(id)
	assert.True(t, len(id) > 0, "비디오 생성 요청이 정상적으로 수행 되지 않았습니다.")
}

func TestCanGetGenerateVideo(t *testing.T) {
	// given : generated_id
	generatedID := ""

	// when : 비디오 전달 요청
	vm := makeModule()
	generatedVideoPath, _ := vm.getGenerateVideo(generatedID)

	// then : 비디오 확인
	assert.True(t, filepath.Base(generatedVideoPath) == "generated_video.mp4", "비디오 파일이 정상적으로 생성되지 않았습니다.")
}

func TestCanMakeReversedVideo(t *testing.T) {
	// given : 유효한 동영상
	// when : 비디오 역재생 영상 요청

	vm := makeModule()
	outputPath, _ := vm.makeReversedVideo()

	// then 비디오 확인
	assert.True(t, filepath.Base(outputPath) == "reversed_video.mp4", "비디오 역 재생 파일이 정상적으로 생성되지 않았습니다.")
}

// func TestCanGenerateVideoContent(t *testing.T) {
// 	vm := makeModule()

// 	vm.GenerateVideoContent()
// }
