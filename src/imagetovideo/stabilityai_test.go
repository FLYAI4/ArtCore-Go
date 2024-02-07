package imagetovideo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

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

// func TestCanPostGenerateVideo(t *testing.T) {
// 	// given : 유효한 토큰 + 이미지
// 	// when : 비디오 생성 요청
// 	vm := makeModule()
// 	id, _ := vm.postGenerateVideo()

// 	// then : 생성 요청 확인(generated_id)
// 	fmt.Println("Generated id : ", id)
// 	assert.True(t, len(id) > 0, "비디오 생성 요청이 정상적으로 수행 되었습니다.")
// }

func TestCanGetGenerateVideo(t *testing.T) {
	// given : generated_id
	// when : 비디오 전달 요청
	token := os.Getenv("STABILITY_AI_TOKEN_KEY")

	generatedVideoPath := "./img/generated_video.mp4"
	generatedID := "789de2faf74ef9b52512966bbcc32b23274fec293adea42f1cf995921aacfc3f"

	var flag int = 202
	url := fmt.Sprintf("https://api.stability.ai/v2alpha/generation/image-to-video/result/%s", generatedID)
	fmt.Println(url)

	for flag == 202 {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Can't create new request.", err)
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "")
		request.Header.Set("authorization", token)

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
			videoContent, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Can't read video content.", err)
			}

			// write video file.
			err = ioutil.WriteFile(generatedVideoPath, videoContent, 0644)
			if err != nil {
				fmt.Println("Can't write video", err)
			}
		case 202:
			fmt.Println("Generation in-progress... automatically try again after 5 sec.")
			time.Sleep(5 * time.Second)
		default:
			fmt.Println("Can't connect api.", err)
			var errorMessage map[string]interface{}
			err = json.NewDecoder(response.Body).Decode(&errorMessage)
			if err != nil {
				fmt.Println("Can't request with Non token", err)
			}
		}
	}

	fmt.Println(flag)
	fmt.Println(generatedVideoPath)

	// then : 비디오 확인
}
