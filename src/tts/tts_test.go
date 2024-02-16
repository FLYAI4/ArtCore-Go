package tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestCanGenerateTts(t *testing.T) {
	mockText := "이 그림은 자연주의와 약간의 인상주의 스타일로 그려진 회화로 보입니다. 이 작품은 전통에서 비롯되었지만, 더 현대적이고 편안한 붓질과 순간을 포착하는 데 관심을 가진 점이 특징입니다. 매체는 캔버스 위에 오일 또는 아크릴 페인트로 보이며, 이는 텍스처와 빛이 표면에서 반사하는 방식에서 나타납니다. 스타일은 약간 단순화되어 있으며, 넓은 색상 영역이 하늘, 초목, 나무의 수직 형태를 구분합니다. 영향을 준 요소로는 인상주의의 측면과 조화와 세부 사항에 대한 더 현대적이고 간소화된 접근법이 포함될 수 있습니다."
	mockFile := "./img/main.mp3"
	os.Remove(mockFile)

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading.")
	}

	token := os.Getenv("OPEN_AI_KEY")

	requestData := map[string]interface{}{
		"model": "tts-1",
		"input": mockText,
		"voice": "nova",
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error to make json data. : ", err)
	}

	// make request
	url := "https://api.openai.com/v1/audio/speech"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Can't make request. : ", err)
	}

	// set request headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	// post request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Can't request to Open AI : ", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		fmt.Println("Can't request to OpenAI. : ", err)
	}

	if response.StatusCode == http.StatusOK {
		audio, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Can't read audo data. : ", err)
		}

		err = os.WriteFile(mockFile, audio, 0644)
		if err != nil {
			fmt.Println("Can't write mp3 file. : ", err)
		}
	}

	// given : 유효한 텍스트

	// when : TTS 생성 요청

	// then : 파일 확인

	_, err = os.Stat(mockFile)
	assert.False(t, os.IsNotExist(err), "gif 파일이 정상적으로 생성되지 않았습니다.")
}
