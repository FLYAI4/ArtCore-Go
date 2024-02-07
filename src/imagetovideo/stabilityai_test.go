package imagetovideo

import (
	"fmt"
	"os"
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
	assert.False(t, os.IsNotExist(err), "리사이즈된 파일이 정상적으로 생성되었습니다.")
}

func TestCanPostGenerateVideo(t *testing.T) {
	// given : 유효한 토큰 + 이미지
	// when : 비디오 생성 요청
	vm := makeModule()
	id, _ := vm.postGenerateVideo()

	// then : 생성 요청 확인(generated_id)
	fmt.Println("Generated id : ", id)
	assert.True(t, len(id) > 0, "비디오 생성 요청이 정상적으로 수행 되었습니다.")
}

func TestCannotRequestGenerateVideoWithNonToken(t *testing.T) {
	// given : 유효하지 않은 토큰 + 이미지

	// when : 비디오 생성 요청

	// then : 토큰 에러 메시지
}

func TestCanGetGenerateVideo(t *testing.T) {
	// given : generated_id

	// when : 비디오 전달 요청

	// then : 비디오 확인
}
