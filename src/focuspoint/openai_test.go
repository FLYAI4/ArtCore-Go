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
	assert.True(t, len(generatedContent.(string)) > 0, "컨텐츠가 정상적으로 생성되지 않았습니다.")
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

func TestCanRefineContent(t *testing.T) {
	// given :

	// when :

	// then :
}

func TestCanGenerateContentAndCoordValue(t *testing.T) {
	// given :

	// when :

	// then :
}
