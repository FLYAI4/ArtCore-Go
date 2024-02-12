package focuspoint

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func makeModule() *FocusPointManager {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading.")
	}

	token := os.Getenv("OPEN_AI_KEY")
	fpm := NewFocusPointManager("./img", token)
	return fpm
}

func TestCanPostGenerateContent(t *testing.T) {
	// given : 유효한 이미지 + 토큰 있는 경우
	fpm := makeModule()

	// when : 생성 요청
	generatedContent, _ := fpm.postGenerateContent()

	// then : 생성된 콘텐츠
	assert.True(t, len(generatedContent) > 0, "컨텐츠가 정상적으로 생성되지 않았습니다.")
	fmt.Println(generatedContent)
}

func TestCannotPostGenerateContentWithNoToken(t *testing.T) {
	// given : 토큰이 없는 경우
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading.")
	}

	token := os.Getenv("OPEN_AI_WRONG_KEY")
	fpm := NewFocusPointManager("./img", token)

	// when : 생성 요청
	generatedContent, err := fpm.postGenerateContent()

	// then : 에러 발생
	assert.True(t, generatedContent == "", "예외가 정상적으로 처리되지 않았습니다.")
	assert.Equal(t, err, fmt.Errorf("non token"), "에러 메시지가 정상적으로 처리되지 않았습니다.")
}

func TestCanRefineMainContent(t *testing.T) {
	// given : 정상적인 content
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
	// when : content 추출 요청
	fpm := makeModule()
	filteredMainContent, _ := fpm.refineMainContent(content)

	// then : content
	fmt.Println(filteredMainContent)
	assert.True(t, len(filteredMainContent) > 0, "컨텐츠가 정상적으로 생성되지 않았습니다.")
}

func TestCanRefindCoordContent(t *testing.T) {
	// given : 정상적인 콘텐츠
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

	fpm := makeModule()
	filteredMainContent, _ := fpm.refineMainContent(content)

	// when : 좌표값 추출 요청
	coordContent, _ := fpm.refineCoordContent(content, filteredMainContent)

	// then : coord content
	fmt.Println(string(coordContent))
	assert.True(t, len(coordContent) > 0, "컨텐츠가 정상적으로 생성되지 않았습니다.")
}

func TestCanGenerateFocusPointContent(t *testing.T) {
	fpm := makeModule()

	fpm.GenerateFocusPointContent()
}
