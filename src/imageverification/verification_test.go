package imageverification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading")
	}
	token := os.Getenv("STABILITY_AI_TOKEN_KEY")
	// given : 유효한 이미지(/storage/{generated_id}/resize.jpg)
	// when : 이미지 upscale

	filePath := filepath.Join("./img", "resized.jpg")
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Can't open the file. : ", err)
	}
	defer file.Close()

	// create new buffer to store the file bytes
	var fileBytes bytes.Buffer
	_, err = fileBytes.ReadFrom(file)
	if err != nil {
		fmt.Println("Can't read file. : ", err)
	}

	// create multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	filePart, err := writer.CreateFormFile("image", "origin.jpg")
	if err != nil {
		fmt.Println("Can't create form file. : ", err)
	}
	_, err = io.Copy(filePart, &fileBytes)
	if err != nil {
		fmt.Println("Can't copy fileBytes to filePart. : ", err)
	}

	changeWeight := 510
	_ = writer.WriteField("prompt", fmt.Sprintf("width: %v", changeWeight))
	_ = writer.WriteField("output_format", "jpeg")
	writer.Close()

	url := "https://api.stability.ai/v2alpha/generation/stable-image/upscale"
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println("Can't make request. : ", err)
	}

	// set header, data
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("authorization", "Bearer "+token)

	// request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Can't request. : ", err)
	}
	defer response.Body.Close()

	fmt.Println("Status Code : ", response.StatusCode)

	// get generatedId
	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		fmt.Println("Can't read response body.", err)
	}

	// decode JSON data into a structure.
	var idStruct struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(buf.Bytes(), &idStruct)
	if err != nil {
		fmt.Println("JSON 디코딩 실패:", err)
	}
	fmt.Println("Generated ID : ", idStruct.ID)

	filePath = filepath.Join("./img", "origin_img.jpg")
	var flag int = 202
	url = fmt.Sprintf("https://api.stability.ai/v2alpha/generation/stable-image/upscale/result/%s", idStruct.ID)
	fmt.Println(url)
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Can't create new request.", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "")
	request.Header.Set("authorization", token)

	for flag == 202 {
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Can't request to get video.", err)
		}
		defer response.Body.Close()

		switch flag = response.StatusCode; flag {
		case 200:
			fmt.Println("Generation finish.")

			// read video data
			imageContent, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Can't read video content.", err)
			}

			// write video file.
			err = os.WriteFile(filePath, imageContent, 0644)
			if err != nil {
				fmt.Println("Can't write video", err)
			}
		case 202:
			fmt.Println("Generation in-progress... automatically try again after 4 sec.")
			time.Sleep(4 * time.Second)
		default:
			//
			fmt.Println(response.Status)
			//
			fmt.Println("Can't connect api.", err)
			var errorMessage map[string]interface{}
			err = json.NewDecoder(response.Body).Decode(&errorMessage)
			if err != nil {
				fmt.Println("Can't request with Non token", err)
			}
		}
	}

	// then : upscale 이미지 생성(/storage/{generated_id}/upscale_img.jpg)
}

func TestCanRetrievalImage(t *testing.T) {
	// given : 유효한 이미지(/storage/{generated_id}/resized.jpg)
	// when : 이미지 유사도 검증
	filePath := filepath.Join("./img", "resized.jpg")
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Can't open the file. : ", err)
	}
	defer file.Close()

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
