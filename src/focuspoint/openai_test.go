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

func TestCannotPostGenerateContentWithNoImage(t *testing.T) {
	// given : 이미지 데이터가 없는 경우

	// when : 생성 요청

	// then : 에러 발생
}

func TestCannotPostGenerateContentWithNoToken(t *testing.T) {
	// given : 토큰이 없는 경우

	// when : 생성 요청

	// then : 에러 발생
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
