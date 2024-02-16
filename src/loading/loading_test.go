package loading

import (
	"fmt"
	img "image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gocv.io/x/gocv"
)

func TestCanSobelMake(t *testing.T) {
	growthRate := float64(0.01)
	background := "black"

	// image load
	image := gocv.IMRead("./img/origin_img.jpg", gocv.IMReadGrayScale)
	if image.Empty() {
		fmt.Println("Can't find image")
	}

	// set sobel config
	cnt := 0
	weightX := float64(0)
	weightY := float64(0)
	delta := float64(0)
	videoFrames := []gocv.Mat{}

	// calculate sobel
	for weightX < 1 || weightY < 1 {
		// x-direction differentiation - vertical mask
		dst1 := gocv.NewMat()
		gocv.Sobel(image, &dst1, gocv.MatTypeCV32F, 1, 0, 3, weightX, delta, 0)

		// y-direction differentiation - horizontal mask
		dst2 := gocv.NewMat()
		gocv.Sobel(image, &dst2, gocv.MatTypeCV32F, 0, 1, 3, weightY, delta, 0)

		// convert absolute value and uint8
		gocv.ConvertScaleAbs(dst1, &dst1, 1, 0)
		gocv.ConvertScaleAbs(dst2, &dst2, 1, 0)

		// combine x-direction, y-direction mask iamge
		mergedImage := gocv.NewMat()
		gocv.AddWeighted(dst1, 0.5, dst2, 0.5, 0, &mergedImage)

		// image count
		cnt++

		// check background
		if background == "white" {
			sobelInverted := gocv.NewMat()
			gocv.BitwiseNot(mergedImage, &sobelInverted)
			videoFrames = append(videoFrames, sobelInverted)
		} else if background == "black" {
			videoFrames = append(videoFrames, mergedImage)
		} else {
			fmt.Println("Wrong background color.")
		}

		// change threshold
		if cnt%2 == 1 {
			weightX += growthRate
		} else {
			weightY += growthRate
		}
	}

	// create gif
	gifPath := "./img/loading.gif"
	os.Remove(gifPath)

	gifFile, err := os.Create(gifPath)
	if err != nil {
		fmt.Println("Error to make gif file. : ", err)
	}
	defer gifFile.Close()

	// generate gif
	gifEncoder := gif.GIF{}

	for _, frame := range videoFrames {
		image, err := frame.ToImage()
		if err != nil {
			fmt.Println("Error to make image : ", err)
		}

		palettedImage := img.NewPaletted(image.Bounds(), color.Palette{color.White, color.Black})
		draw.Draw(palettedImage, palettedImage.Rect, image, image.Bounds().Min, draw.Src)

		gifEncoder.Delay = append(gifEncoder.Delay, 0)
		gifEncoder.Image = append(gifEncoder.Image, palettedImage)
	}

	err = gif.EncodeAll(gifFile, &gifEncoder)
	if err != nil {
		fmt.Println("Error to encdoe image : ", err)
	}
	fmt.Println("Create")

	_, err = os.Stat(gifPath)
	assert.False(t, os.IsNotExist(err), "gif 파일이 정상적으로 생성되지 않았습니다.")
	// given : 유효한 이미지

	// when : loading 이미지 생성 요청

	// then : 생성 완료
}
