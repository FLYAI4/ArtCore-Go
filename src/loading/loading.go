package loading

import (
	"fmt"
	img "image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"path/filepath"
	"sync"

	"github.com/robert-min/ArtCore-Go/src/pb"
	"gocv.io/x/gocv"
)

type LoadingManager struct {
	userFolderPath string
	growthRate     float64
	background     string
}

func NewLoadingManager(userFolderPath string) *LoadingManager {
	return &LoadingManager{
		userFolderPath: userFolderPath,
		growthRate:     float64(0.015),
		background:     "black",
	}
}

func (lm *LoadingManager) GetLodingGif(wg *sync.WaitGroup, stream pb.StreamService_GeneratedContentStreamServer) error {
	videoFrames, err := lm.generateSobelFrame()
	if err != nil {
		fmt.Println("Generate sobel frame error : ", err)
		return err
	}

	gifFilePath, err := lm.generateGif(videoFrames)
	if err != nil {
		fmt.Println("Genreate gif error : ", err)
		return err
	}

	gifBytes, err := os.ReadFile(gifFilePath)
	if err != nil {
		fmt.Println("Failed to read gif error: ", err)
	}

	// send: gif content
	if err := stream.Send(&pb.Response{Tag: "gif", Data: gifBytes}); err != nil {
		fmt.Println("Failed to send response: ", err)
	}

	fmt.Println("Finsish gRPC")
	wg.Done()
	return nil
}

func (lm *LoadingManager) generateSobelFrame() ([]gocv.Mat, error) {
	// image load
	filePath := filepath.Join(lm.userFolderPath, "origin_img.jpg")

	image := gocv.IMRead(filePath, gocv.IMReadGrayScale)
	if image.Empty() {
		fmt.Println("Can't find image")
		return nil, fmt.Errorf("error to load image")
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
		if lm.background == "white" {
			sobelInverted := gocv.NewMat()
			gocv.BitwiseNot(mergedImage, &sobelInverted)
			videoFrames = append(videoFrames, sobelInverted)
		} else if lm.background == "black" {
			videoFrames = append(videoFrames, mergedImage)
		} else {
			fmt.Println("Wrong background color.")
			return nil, fmt.Errorf("error to wrong color")
		}

		// change threshold
		if cnt%2 == 1 {
			weightX += lm.growthRate
		} else {
			weightY += lm.growthRate
		}
	}
	return videoFrames, nil
}

func (lm *LoadingManager) generateGif(videoFrames []gocv.Mat) (string, error) {
	// create gif
	filePath := filepath.Join(lm.userFolderPath, "loading.gif")
	os.Remove(filePath)

	gifFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error to make gif file. : ", err)
		return "", err
	}
	defer gifFile.Close()

	// generate gif
	gifEncoder := gif.GIF{}

	for _, frame := range videoFrames {
		image, err := frame.ToImage()
		if err != nil {
			fmt.Println("Error to make image : ", err)
			return "", err
		}

		palettedImage := img.NewPaletted(image.Bounds(), color.Palette{color.White, color.Black})
		draw.Draw(palettedImage, palettedImage.Rect, image, image.Bounds().Min, draw.Src)

		gifEncoder.Delay = append(gifEncoder.Delay, 0)
		gifEncoder.Image = append(gifEncoder.Image, palettedImage)
	}

	err = gif.EncodeAll(gifFile, &gifEncoder)
	if err != nil {
		fmt.Println("Error to encdoe image : ", err)
		return "", err
	}
	fmt.Println("Create")
	return filePath, nil
}
