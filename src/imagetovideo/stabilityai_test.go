package imagetovideo

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"testing"

	"github.com/nfnt/resize"
	"github.com/stretchr/testify/assert"
)

func TestCanResizeImage(t *testing.T) {
	// given : 유효한 이미지
	// when : 이미지 사이즈 변경 요청
	// open local file
	file, err := os.Open("./img/origin_img.jpg")
	if err != nil {
		fmt.Println("Can't open the file.", err)
		return
	}
	defer file.Close()

	// image decoding
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Can't decoding the image.", err)
		return
	}

	// image resizing
	changeWidth := 768
	changeHeight := 768
	resizeImg := resize.Resize(uint(changeWidth), uint(changeHeight), img, resize.Lanczos3)

	// create new file
	newFile, err := os.Create("./img/resized_img.jpg")
	if err != nil {
		fmt.Println("Can't create new file.", err)
		return
	}
	defer newFile.Close()

	// save the image to file
	err = jpeg.Encode(newFile, resizeImg, nil)
	if err != nil {
		fmt.Println("Can't save the image to file", err)
		return
	}
	fmt.Println("The image has been successfully resized and saved.")

	// then : 이미지 저장 파일 확인
	resizedFilePath := "./img/resized_img.jpg"
	_, err = os.Stat(resizedFilePath)
	if os.IsNotExist(err) {
		assert.Fail(t, "리사이즈된 파일이 존재하지 않습니다.")
	} else if err != nil {
		assert.Fail(t, "파일 상태를 확인하는 중 오류가 발생했습니다: "+err.Error())
	} else {
		assert.True(t, true, "리사이즈된 파일이 정상적으로 생성되었습니다.")
	}
}

func TestCanRequestGenerateVideo(t *testing.T) {
	// given : 유효한 토큰 + 이미지

	// when : 비디오 생성 요청

	// then : 생성 요청 확인(generated_id)
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
