package imagetovideo

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

type VideoManager struct {
	changeWidth    int
	changeHeight   int
	userFolderPath string
}

func NewVideoManager(userFolderPath string) *VideoManager {
	return &VideoManager{
		changeWidth:    768,
		changeHeight:   768,
		userFolderPath: userFolderPath,
	}
}

func (vm *VideoManager) GenerateVideoContent() {
	// resize image to 768 * 768
	_, err := vm.resizeImage()
	if err != nil {
		fmt.Println("Resize image error.", err)
	}
}

func (vm *VideoManager) resizeImage() (string, error) {
	filePath := filepath.Join(vm.userFolderPath, "origin_img.jpg")

	// open local file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Can't open the file.", err)
		return "", err
	}
	defer file.Close()

	// image decoding
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Can't decoding the image.", err)
		return "", err
	}

	// image resizing
	resizeImg := resize.Resize(uint(vm.changeHeight), uint(vm.changeHeight), img, resize.Lanczos3)

	// create new file
	filePath = filepath.Join(vm.userFolderPath, "resized_img.jpg")
	newFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Can't create new file.", err)
		return "", err
	}
	defer newFile.Close()

	// save the image to file
	err = jpeg.Encode(newFile, resizeImg, nil)
	if err != nil {
		fmt.Println("Can't save the image to file", err)
		return "", err
	}
	fmt.Println("The image has been successfully resized and saved.")

	return filePath, err
}
