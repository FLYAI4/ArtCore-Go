package imageverification

import "testing"

func TestCanResizeImageWithJPG(t *testing.T) {
	// given : 유효한 이미지(/storage/{generated_id}/origin.jpg)

	// when : 이미지 리사이즈

	// then : 리사이즈 이미지 생성(/storage/{generated_id}/resize.jpg)
}

func TestCanUpscaleImage(t *testing.T) {
	// given : 유효한 이미지(/storage/{generated_id}/resize.jpg)

	// when : 이미지 upscale

	// then : upscale 이미지 생성(/storage/{generated_id}/upscale_img.jpg)
}

func TestCanRetrievalImage(t *testing.T) {
	// given : 유효한 이미지(/storage/{generated_id}/upscale_img.jpg)

	// when : 이미지 유사도 검증

	// then : 유사도 수치 확인
}

func TestCanVerificationImage(t *testing.T) {
	// given : 유효한 이미지

	// when : 이미지 검증

	// then : 정상 응답
}

func TestCannotVerificationImage(t *testing.T) {
	// given : 유효한 이미지 + threshold 아래 값

	// when : 이미지 검증

	// then : 이미지 재촬영 요청 + 폴더 삭제
}
