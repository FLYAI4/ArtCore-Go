package focuspoint

import "testing"

func TestCannotGetGenerateContentWithNoImage(t *testing.T) {
	// given : 이미지 데이터가 없는 경우

	// when : 생성 요청

	// then : 에러 발생
}

func TestCannotGetGenerateContentWithNoToken(t *testing.T) {
	// given : 토큰이 없는 경우

	// when : 생성 요청

	// then : 에러 발생
}

func TestCanGetGenerateContent(t *testing.T) {
	// given : 유효한 이미지 + 토큰 있는 경우

	// when : 생성 요청

	// then : 생성된 콘텐츠
}

func TestCanRefineContent(t *testing.T) {
	// given : 

	// when : 

	// then : 
}

func TestCanGenerateContentAndCoordValue(t *testing.T)) {
	// given : 

	// when :

	// then : 
}