package imagetovideo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanResizeImage(t *testing.T) {
	// given : 유효한 이미지
	// when : 이미지 사이즈 변경 요청
	vm := NewVideoManager("./img")
	resizedFilePath, _ := vm.resizeImage()

	// then : 이미지 저장 파일 확인
	_, err := os.Stat(resizedFilePath)
	assert.False(t, os.IsNotExist(err), "리사이즈된 파일이 정상적으로 생성되었습니다.")
}

func TestCanRequestGenerateVideo(t *testing.T) {
	// given : 유효한 토큰 + 이미지
	// when : 비디오 생성 요청
	resizedImagePath := "./img/resized_img.jpg"

	// open file
	file, err := os.Open(resizedImagePath)
	if err != nil {
		fmt.Println("Can't open resized file.", err)
	}
	defer file.Close()

	// get fileInfo
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Can't get file info.", err)
	}

	// read fileInfo
	fileBytes := make([]byte, fileInfo.Size())
	_, err = file.Read(fileBytes)
	if err != nil {
		fmt.Println("Can't read file info.", err)
	}

	// put file data into buffer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	formFile, err := writer.CreateFormFile("image", fileInfo.Name())
	if err != nil {
		fmt.Println("Can't create form file.", err)
	}
	_, err = io.Copy(formFile, bytes.NewReader(fileBytes))
	if err != nil {
		fmt.Println("Can't put file data into buffer.", err)
	}
	writer.Close()

	// create request
	token := ""
	url := "https://api.stability.ai/v2alpha/generation/image-to-video"
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println("Can't create new request.", err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("authorization", "Bearer "+token)

	// post request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Can't request with wrong request", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Success to request.")
	} else if response.StatusCode >= 400 && response.StatusCode < 500 {
		var errorMessage map[string]interface{}
		err = json.NewDecoder(response.Body).Decode(&errorMessage)
		if err != nil {
			fmt.Println("Can't request with Non token", err)
		}
	} else {
		var errorMessage map[string]interface{}
		err = json.NewDecoder(response.Body).Decode(&errorMessage)
		if err != nil {
			fmt.Println("Can't request with unknown error", err)
		}
	}

	// then : 생성 요청 확인(generated_id)
	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		fmt.Println("Can't read response body.", err)
	}
	fmt.Println("응답 값 : ", buf.String())

	assert.Equal(t, response.StatusCode, http.StatusOK, "잘못된 Status code.")
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
