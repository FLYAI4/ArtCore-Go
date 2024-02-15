package focuspoint

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/robert-min/ArtCore-Go/src/pb"
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

func (fpm *FocusPointManager) GenerateFocusPointContent(wg *sync.WaitGroup, stream pb.StreamService_GeneratedContentStreamServer) {
	content, err := fpm.postGenerateContent()
	if err != nil {
		fmt.Println("Post generateContent error: ", err)
	}

	mainContent, err := fpm.refineMainContent(content)
	if err != nil {
		fmt.Println("Refine mainContent error: ", err)
	}
	// send: mainContent
	if err := stream.Send(&pb.Response{Tag: "content", Data: []byte(mainContent)}); err != nil {
		fmt.Println("Failed to send response: ", err)
	}

	coorContent, err := fpm.refineCoordContent(content, mainContent)
	if err != nil {
		fmt.Println("Refine coordContent error: ", err)
	}
	// send: coordContent
	if err := stream.Send(&pb.Response{Tag: "coord", Data: coorContent}); err != nil {
		fmt.Println("Failed to send response: ", err)
	}
	wg.Done()
}

func (fpm *FocusPointManager) postGenerateContent() (string, error) {
	// open file
	filePath := filepath.Join(fpm.userFolderPath, "origin_img.jpg")
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Can't open the file.", err)
		return "", err
	}
	defer file.Close()

	// read file
	imageData, err := io.ReadAll(file)
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
	} else if response.StatusCode == 429 {
		fmt.Println("Can't request to OpenAI with non token.")
		return "", fmt.Errorf("non token")
	}

	responseData, err := io.ReadAll(response.Body)
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
	messageContent := ""
	if choices, ok := data["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					messageContent = content
				}
			}
		}
	}

	return messageContent, nil
}

func (fpm *FocusPointManager) refineMainContent(content string) (string, error) {
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
	// TODO : filteredMainContent 전송(byte 타입으로)
	return filteredMainContent, nil
}

func (fpm *FocusPointManager) refineCoordContent(content string, filteredMainContent string) ([]byte, error) {
	// find start
	startIndex := strings.Index(content, "```json")
	if startIndex == -1 {
		fmt.Println("Start index not found.")
		return []byte("{}"), fmt.Errorf("no json format error")
	}
	startIndex += 7

	// find end
	endIndex := strings.Index(content[startIndex:], "```")
	if endIndex == -1 {
		fmt.Println("End index not found")
		return []byte("{}"), fmt.Errorf("no json format error")
	}

	// decode json
	jsonStr := content[startIndex : startIndex+endIndex]
	jsonStr = strings.ReplaceAll(jsonStr, "\n", "")
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return []byte("{}"), err
	}

	// make coord content
	changeJsonData := make(map[string]interface{})
	for key, value := range jsonData {
		changeJsonData[key] = map[string]interface{}{
			"content": findCoordContent(key, filteredMainContent),
			"coord":   value,
		}
	}
	changedJSONBytes, err := json.Marshal(changeJsonData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return []byte("{}"), err
	}
	return changedJSONBytes, nil

}

func findCoordContent(key string, filteredMainContent string) string {
	var coordSentences []string

	parts := strings.Split(key, "_")
	lastPart := parts[len(parts)-1]

	sentences := strings.Split(filteredMainContent, ".")
	for _, sentence := range sentences {
		if strings.Contains(sentence, lastPart) {
			coordSentences = append(coordSentences, sentence)
		}
	}

	result := strings.Join(coordSentences, ". ")
	return result
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
