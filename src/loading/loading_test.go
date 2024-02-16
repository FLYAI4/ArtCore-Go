package loading

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanSobelMake(t *testing.T) {
	// given : 유효한 이미지
	lm := NewLoadingManager("./img")

	// when : loading 이미지 생성 요청
	lm.GetLodingGif()

	// then : 생성 완료
	_, err := os.Stat("./img/loading.gif")
	assert.False(t, os.IsNotExist(err), "gif 파일이 정상적으로 생성되지 않았습니다.")
}
