package imagetovideo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

type VideoManager struct {
	changeWidth    int
	changeHeight   int
	userFolderPath string
	token          string
}

func NewVideoManager(userFolderPath string, token string) *VideoManager {
	return &VideoManager{
		changeWidth:    768,
		changeHeight:   768,
		userFolderPath: userFolderPath,
		token:          token,
	}
}

func (vm *VideoManager) GenerateVideoContent() {
	// resize image to 768 * 768
	_, err := vm.resizeImage()
	if err != nil {
		fmt.Println("Resize image error: ", err)
	}
	_, err = vm.postGenerateVideo()
	if err != nil {
		fmt.Println("Post generate video error: ", err)
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

	return filePath, nil
}

func (vm *VideoManager) postGenerateVideo() (string, error) {
	filePath := filepath.Join(vm.userFolderPath, "resized_img.jpg")

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Can't open resized file.", err)
		return "", err
	}
	defer file.Close()

	// get fileInfo
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Can't get file info.", err)
		return "", err
	}

	// read fileInfo
	fileBytes := make([]byte, fileInfo.Size())
	_, err = file.Read(fileBytes)
	if err != nil {
		fmt.Println("Can't read file info.", err)
		return "", err
	}

	// put file data into buffer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	formFile, err := writer.CreateFormFile("image", fileInfo.Name())
	if err != nil {
		fmt.Println("Can't create form file.", err)
		return "", err
	}
	_, err = io.Copy(formFile, bytes.NewReader(fileBytes))
	if err != nil {
		fmt.Println("Can't put file data into buffer.", err)
		return "", err
	}
	writer.Close()

	// create request
	url := "https://api.stability.ai/v2alpha/generation/image-to-video"
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println("Can't create new request.", err)
		return "", err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("authorization", "Bearer "+vm.token)

	// post request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Can't request with wrong request", err)
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Success to request.")
	} else if response.StatusCode >= 400 && response.StatusCode < 500 {
		var errorMessage map[string]interface{}
		err = json.NewDecoder(response.Body).Decode(&errorMessage)
		if err != nil {
			fmt.Println("Can't request with Non token", err)
			return "", err
		}
	} else {
		var errorMessage map[string]interface{}
		err = json.NewDecoder(response.Body).Decode(&errorMessage)
		if err != nil {
			fmt.Println("Can't request with unknown error", err)
			return "", err
		}
	}

	// get generatedId
	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		fmt.Println("Can't read response body.", err)
	}

	return buf.String(), nil
}
