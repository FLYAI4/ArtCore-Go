package imageverification

import (
	"fmt"
	"image"

	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/nfnt/resize"
	"github.com/stretchr/testify/assert"
)

func TestCanResizeImageWithJPG(t *testing.T) {
	// given : 유효한 이미지(/storage/{generated_id}/origin.jpg)
	// when : 이미지 리사이즈

	// file open
	filePath := filepath.Join("./img", "origin.jpg")

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Can't open the file. : ", err)
	}
	defer file.Close()

	// image decoding
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Can't decoding the image. : ", err)
	}

	// image resizing
	changeWeight := 510
	changeHeight := 680
	resizeImg := resize.Resize(uint(changeWeight), uint(changeHeight), img, resize.Lanczos3)

	// create new file
	filePath = filepath.Join("./img", "resized.jpg")
	newFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Can't create new file. : ", err)
	}
	defer newFile.Close()

	// save the image to file
	err = jpeg.Encode(newFile, resizeImg, nil)
	if err != nil {
		fmt.Println("Can't save the image to file. : ", err)
	}

	// then : 리사이즈 이미지 생성(/storage/{generated_id}/resize.jpg)
	filePath = filepath.Join("./img", "resized.jpg")
	file, err = os.Open(filePath)
	if err != nil {
		fmt.Println("Can't open the file. : ", err)
	}

	img, _, err = image.Decode(file)
	if err != nil {
		fmt.Println("Can't decoding the image. : ", err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	fmt.Printf("이미지 크기: width = %d, height = %d\n", width, height)
	assert.True(t, width == changeWeight, "크기가 정상적으로 변경되지 않았습니다.")
	assert.True(t, height == changeHeight, "크기가 정상적으로 변경되지 않았습니다.")
}

func TestCanUpscaleImage(t *testing.T) {
	// given : 유효한 이미지(/storage/{generated_id}/resize.jpg)

	// when : 이미지 upscale

	// then : upscale 이미지 생성(/storage/{generated_id}/upscale_img.jpg)
}

func TestCanRetrievalImage(t *testing.T) {
	// given : 유효한 이미지(/storage/{generated_id}/upscale_img.jpg)

	// when : 이미지 유사도 검증

	// then : 유사도 수치 확인
}

func TestCanVerificationImage(t *testing.T) {
	// given : 유효한 이미지

	// when : 이미지 검증

	// then : 정상 응답
}

func TestCannotVerificationImage(t *testing.T) {
	// given : 유효한 이미지 + threshold 아래 값

	// when : 이미지 검증

	// then : 이미지 재촬영 요청 + 폴더 삭제
}
