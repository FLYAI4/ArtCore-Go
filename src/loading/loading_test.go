package loading

import (
	"os"
	"sync"
	"testing"

	"github.com/robert-min/ArtCore-Go/src/pb"
	"github.com/stretchr/testify/assert"
)

func TestCanSobelMake(t *testing.T) {
	// given : 유효한 이미지
	lm := NewLoadingManager("./img")
	var wg sync.WaitGroup
	var stream pb.StreamService_GeneratedContentStreamServer
	wg.Add(1)

	// when : loading 이미지 생성 요청
	lm.GetLodingGif(&wg, stream)
	wg.Wait()

	// then : 생성 완료
	_, err := os.Stat("./img/loading.gif")
	assert.False(t, os.IsNotExist(err), "gif 파일이 정상적으로 생성되지 않았습니다.")
}
