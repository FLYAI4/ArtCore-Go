package focuspoint

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FocusPointManager struct {
	userFolderPath string
	token          string
}

func NewFocusPointManager(userFolderPath string, token string) *FocusPointManager {
	return &FocusPointManager{
		userFolderPath: userFolderPath,
		token:          token,
	}
}

func (fpm *FocusPointManager) postGenerateContent() (interface{}, error) {
	// open file
	filePath := filepath.Join(fpm.userFolderPath, "origin_img.jpg")
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Can't open the file.", err)
		return "", err
	}
	defer file.Close()

	// read file
	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Can't read the file.", err)
		return "", err
	}

	// encode Base64
	encodedImage := base64.StdEncoding.EncodeToString(imageData)

	// make json data
	requestData := map[string]interface{}{
		"model": "gpt-4-vision-preview",
		"messages": []interface{}{
			map[string]interface{}{
				"role": "user",
				"content": []interface{}{
					map[string]interface{}{
						"type": "text",
						"text": makePrompt(),
					},
					map[string]interface{}{
						"type": "image_url",
						"image_url": map[string]interface{}{
							"url": fmt.Sprintf("data:image/jpeg;base64,%s", encodedImage),
						},
					},
				},
			},
		},
		"max_tokens": 800,
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Can't make json data.")
		return "", err
	}

	// make request
	url := "https://api.openai.com/v1/chat/completions"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Can't make request")
		return "", err
	}

	// set request headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+fpm.token)

	// post request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Can't request to OpenAI", err)
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode == 400 {
		fmt.Println("Can't request to OpenAI with non image.")
		return "", fmt.Errorf("non image")
	} else if response.StatusCode == 404 {
		fmt.Println("Can't request to OpenAI with non token.")
		return "", fmt.Errorf("non token")
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Can't read reponse body.", err)
		return "", err
	}
	// parse json data
	var data map[string]interface{}
	if err := json.Unmarshal(responseData, &data); err != nil {
		fmt.Println("Can't parse json data.", err)
		return "", err
	}
	// ['choices'][0]["message"]["content"]
	messageContent := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"]
	return messageContent, nil
}

func (fpm *FocusPointManager) refineContent(content string) (string, error) {
	// find main content
	splitContent := strings.Split(content, ":")
	if len(splitContent) == 0 {
		fmt.Println("Can't split data")
		return "", fmt.Errorf("non split content")
	}
	mainContent := strings.TrimSpace(splitContent[0])

	// filter words
	words := []string{"cannot", "AI", "do not", "can't", "json", "JSON", "{", "Unfortunately", "coordinates", "However", "keyword", "keywords"}
	filteredMainContent := filterSentences(mainContent, words)
	// TODO : filteredMainContent 전송
	return filteredMainContent, nil
}

func filterSentences(content string, words []string) string {
	var filteredSentences []string

	sentences := strings.Split(content, ".")
OuterLoop:
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		for _, word := range words {
			if strings.Contains(sentence, word) {
				continue OuterLoop
			}
		}
		filteredSentences = append(filteredSentences, sentence)
	}

	result := strings.Join(filteredSentences, ". ")
	return result
}

func makePrompt() string {
	promtText := []string{
		"You are an expert art historian with vast knowledge about artists throughout history who revolutionized their craft.",
		"You will begin by briefly summarizing the personal life and achievements of the artist.",
		"Then you will go on to explain the medium, style, and influences of their works.",
		"Then you will provide short descriptions of what they depict and any notable characteristics they might have.",
		"Fianlly identify THREE keywords in the picture and provide each coordinate of the keywords in the last sentence.",
		"For example, Give the coordinate value of the keywords in json format.",
		"if the keyword is pretty_woman and big_ball, value is  ```json{\"pretty_woman\", [[x0,y0,x1,y1]], \"big_ball\", [[x0,y0,x1,y1], [x2,y2,x3,y3]]}```",
		"The values ​​entered in x0, y0, x1, y1 are unconditionally the coordinate values ​​of each keyword.",
	}
	return strings.Join(promtText, " ")
}
