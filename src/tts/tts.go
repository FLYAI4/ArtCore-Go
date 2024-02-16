package tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type AudioManager struct {
	userFolderPath string
	token          string
}

func NewAudioManager(userFolderPath string, token string) *AudioManager {
	return &AudioManager{
		userFolderPath: userFolderPath,
		token:          token,
	}
}

func (am *AudioManager) GetAudioContent(text string) error {
	_, err := am.generateTts(text)
	if err != nil {
		fmt.Println("Error to generate tts. : ", err)
		return err
	}
	return nil
}

func (am *AudioManager) generateTts(text string) (string, error) {
	filePath := filepath.Join(am.userFolderPath, "main.mp3")

	requestData := map[string]interface{}{
		"model": "tts-1",
		"input": text,
		"voice": "nova",
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error to make json data. : ", err)
		return "", err
	}

	// make request
	url := "https://api.openai.com/v1/audio/speech"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Can't make request. : ", err)
		return "", err
	}

	// set request headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+am.token)

	// post request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Can't request to Open AI : ", err)
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		fmt.Println("Can't request to OpenAI. : ", err)
		return "", err
	}

	if response.StatusCode == http.StatusOK {
		// save the file
		audio, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Can't read audo data. : ", err)
			return "", err
		}

		err = os.WriteFile(filePath, audio, 0644)
		if err != nil {
			fmt.Println("Can't write mp3 file. : ", err)
			return "", err
		}
	}
	return filePath, nil
}
