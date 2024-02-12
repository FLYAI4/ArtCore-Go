package focuspoint

import (
	"fmt"
	"strings"
	"testing"
)

// func makeModule() *FocusPointManager {
// 	err := godotenv.Load()
// 	if err != nil {
// 		fmt.Println("Error loading.")
// 	}

// 	token := os.Getenv("OPEN_AI_KEY")
// 	fpm := NewFocusPointManager("./img", token)
// 	return fpm
// }

// func TestCanPostGenerateContent(t *testing.T) {
// 	// given : 유효한 이미지 + 토큰 있는 경우
// 	fpm := makeModule()

// 	// when : 생성 요청
// 	generatedContent, _ := fpm.postGenerateContent()

// 	// then : 생성된 콘텐츠
// 	assert.True(t, len(generatedContent.(string)) > 0, "컨텐츠가 정상적으로 생성되지 않았습니다.")
// 	fmt.Println(generatedContent)
// }

// func TestCannotPostGenerateContentWithNoToken(t *testing.T) {
// 	// given : 토큰이 없는 경우
// 	err := godotenv.Load()
// 	if err != nil {
// 		fmt.Println("Error loading.")
// 	}

// 	token := os.Getenv("OPEN_AI_WRONG_KEY")
// 	fpm := NewFocusPointManager("./img", token)

// 	// when : 생성 요청
// 	generatedContent, err := fpm.postGenerateContent()

// 	// then : 에러 발생
// 	assert.True(t, generatedContent == "", "예외가 정상적으로 처리되지 않았습니다.")
// 	assert.Equal(t, err, fmt.Errorf("non token"), "에러 메시지가 정상적으로 처리되지 않았습니다.")

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

func TestCanRefineMainContent(t *testing.T) {
	content := "As an AI developed by OpenAI, I don't have access to external databases to identify specific artwork or artists, and can't comment on the personal life or achievements of this artist without that information. However, I can analyze the image based on art theory and the visible elements.\n" +
		"This piece of art seems to be a vibrant landscape painting. The medium appears to be acrylic or oil paint on canvas, evident from the texture and the way the light reflects off the surface. The style is somewhat stylized with a slight hint of naïve art due to the simplistic representation of the scenery and the bright, almost uniform colors.\n" +
		"Influences in the work might include Fauvism for the bold color choices and Asian landscape painting for the terraced fields, which are typically seen in East Asian countries like China and Vietnam. Such fields are often used to grow rice and are sculpted into hilly or mountainous terrain.\n" +
		"The painting depicts a lush, green landscape with staggered rows of terraced fields that are indicative of careful human cultivation. A winding river or lake weaves through the mountains into the distance, bringing a sense of depth and tranquility. Hot air balloons float serenely in the sky, contributing to the overall peaceful and picturesque scenery. Notably, the scene has a harmonious palette, primarily composed of greens and blues, which evokes a calm and natural setting.\n" +
		"For the keywords, let's select \"hot_air_balloons,\" \"terraced_fields,\" and \"blue_water.\" Please note that I am not able to provide actual pixel coordinates, so the following values will be illustrative approximations based on the image:\n" +
		"\n" +
		"```json\n" +
		"{\n" +
		"\"hot_air_balloons\": [[100,40,180,100], [200,30,270,85], [320,20,370,70]],\n" +
		"\"terraced_fields\": [[0,200,800,600]],\n" +
		"\"blue_water\": [[150,300,650,400]]\n" +
		"}\n" +
		"```\n"

	// find main content
	splitContent := strings.Split(content, ":")
	mainContent := strings.TrimSpace(splitContent[0])
	// if len(splitContent) > 0 {
	// 	mainContent := strings.TrimSpace(splitContent[0])
	// 	fmt.Println(mainContent)
	// }

	// filter words
	words := []string{"cannot", "AI", "do not", "can't", "json", "JSON", "{", "Unfortunately", "coordinates", "However", "keyword", "keywords"}
	filteredContent := filterSentences(mainContent, words)
	fmt.Println(filteredContent)

	// given : 정상적인 content

	// when : content 추출 요청

	// then : content
}

// func TestCanRefindCoordValue(t *testing.T) {
// 	content := "As an AI developed by OpenAI, I don't have access to external databases to identify specific artwork or artists, and can't comment on the personal life or achievements of this artist without that information. However, I can analyze the image based on art theory and the visible elements.\n" +
// 		"This piece of art seems to be a vibrant landscape painting. The medium appears to be acrylic or oil paint on canvas, evident from the texture and the way the light reflects off the surface. The style is somewhat stylized with a slight hint of naïve art due to the simplistic representation of the scenery and the bright, almost uniform colors.\n" +
// 		"Influences in the work might include Fauvism for the bold color choices and Asian landscape painting for the terraced fields, which are typically seen in East Asian countries like China and Vietnam. Such fields are often used to grow rice and are sculpted into hilly or mountainous terrain.\n" +
// 		"The painting depicts a lush, green landscape with staggered rows of terraced fields that are indicative of careful human cultivation. A winding river or lake weaves through the mountains into the distance, bringing a sense of depth and tranquility. Hot air balloons float serenely in the sky, contributing to the overall peaceful and picturesque scenery. Notably, the scene has a harmonious palette, primarily composed of greens and blues, which evokes a calm and natural setting.\n" +
// 		"For the keywords, let's select \"hot_air_balloons,\" \"terraced_fields,\" and \"blue_water.\" Please note that I am not able to provide actual pixel coordinates, so the following values will be illustrative approximations based on the image:\n" +
// 		"\n" +
// 		"```json\n" +
// 		"{\n" +
// 		"\"hot_air_balloons\": [[100,40,180,100], [200,30,270,85], [320,20,370,70]],\n" +
// 		"\"terraced_fields\": [[0,200,800,600]],\n" +
// 		"\"blue_water\": [[150,300,650,400]]\n" +
// 		"}\n" +
// 		"```\n"

// 	// find main content
// 	splitContent := strings.Split(content, ":")
// 	mainContent := strings.TrimSpace(splitContent[0])

// 	// find start
// 	startIndex := strings.Index(content, "```json")
// 	if startIndex == -1 {
// 		fmt.Println("Start index not found.")
// 	}
// 	startIndex += 7

// 	// find end
// 	endIndex := strings.Index(content[startIndex:], "```")
// 	if endIndex == -1 {
// 		fmt.Println("End index not found")
// 	}

// 	// decode json
// 	jsonStr := content[startIndex : startIndex+endIndex]
// 	jsonStr = strings.ReplaceAll(jsonStr, "\n", "")
// 	var jsonData map[string]interface{}
// 	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
// 		fmt.Println("Error decoding JSON:", err)
// 		return
// 	}

// 	fmt.Println("==================")
// 	sentences := strings.Split(mainContent, ".")
// 	for _, sentence := range sentences {
// 		for key := range jsonData {
// 			if strings.Contains(sentence, key) {
// 				fmt.Println(strings.TrimSpace(sentence))
// 			}
// 		}
// 	}

// 	// given : 정상적인 콘텐츠

// 	// when : 좌표값 추출 요청

// 	// then : coord content
// }
