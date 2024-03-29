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
	"sync"
	"time"

	"github.com/nfnt/resize"
	"github.com/robert-min/ArtCore-Go/src/pb"
)

type VideoManager struct {
	changeWidth    int
	changeHeight   int
	userFolderPath string
	token          string
}

// NewVideoManager make VideomManager struct.
func NewVideoManager(userFolderPath string, token string) *VideoManager {
	return &VideoManager{
		changeWidth:    768,
		changeHeight:   768,
		userFolderPath: userFolderPath,
		token:          token,
	}
}

// Generate video content with StabilityAI and openCV.
func (vm *VideoManager) GenerateVideoContent(wg *sync.WaitGroup, stream pb.StreamService_GeneratedContentStreamServer) error {
	// resize image to 768 * 768
	_, err := vm.resizeImage()
	if err != nil {
		fmt.Println("Resize image error: ", err)
		return err
	}
	// post generated video to stability AI
	id, err := vm.postGenerateVideo()
	if err != nil {
		fmt.Println("Post generate video error: ", err)
		return err
	}

	// TODO : Front에서 응답을 종료할 수 있는지 확인!!
	// get generated video from stability AI
	_, err = vm.getGenerateVideo(id)
	if err != nil {
		fmt.Println("Get generate video error: ", err)
	}

	// make long reversed video to openCV
	outFilePath, err := vm.makeReversedVideo()
	if err != nil {
		fmt.Println("Make reversed video error: ", err)
	}

	// TODO : send video 별도 api로 수정
	// read video to byte
	videoBytes, err := os.ReadFile(outFilePath)
	if err != nil {
		fmt.Println("Failed to read video error: ", err)
	}

	// send: video content
	if err := stream.Send(&pb.Response{Tag: "video", Data: videoBytes}); err != nil {
		fmt.Println("Failed to send response: ", err)
	}

	// send: finish
	if err := stream.Send(&pb.Response{Tag: "finish", Data: []byte("finished")}); err != nil {
		fmt.Println("Failed to send response: ", err)
	}

	wg.Done()
	return nil
}

// func (vm *VideoManager) GetVideoContent(stream pb.StreamService_VideoContentStreamServer) {
// 	// make long reversed video to openCV
// 	outFilePath, err := vm.makeReversedVideo()
// 	if err != nil {
// 		fmt.Println("Make reversed video error: ", err)
// 	}

// 	// TODO : send video 별도 api로 수정
// 	// read video to byte
// 	videoBytes, err := os.ReadFile(outFilePath)
// 	if err != nil {
// 		fmt.Println("Failed to read video error: ", err)
// 	}

// 	// send: video content
// 	if err := stream.Send(&pb.Response{Tag: "video", Data: videoBytes}); err != nil {
// 		fmt.Println("Failed to send response: ", err)
// 	}
// }

// Resize image to fit (768, 768)
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

// Post generate video with stabilityAI
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

	// decode JSON data into a structure.
	var idStruct struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(buf.Bytes(), &idStruct)
	if err != nil {
		fmt.Println("JSON 디코딩 실패:", err)
		return "", err
	}

	return idStruct.ID, nil
}

// Get generate video with stabilityAI
func (vm *VideoManager) getGenerateVideo(generatedID string) (string, error) {
	filePath := filepath.Join(vm.userFolderPath, "generated_video.mp4")
	deleteFile(filePath)

	var flag int = 202
	url := fmt.Sprintf("https://api.stability.ai/v2alpha/generation/image-to-video/result/%s", generatedID)
	fmt.Println(url)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Can't create new request.", err)
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "")
	request.Header.Set("authorization", vm.token)

	for flag == 202 {
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Can't request to get video.", err)
			return "", err
		}
		defer response.Body.Close()

		switch flag = response.StatusCode; flag {
		case 200:
			fmt.Println("Generation finish.")

			// read video data
			videoContent, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Can't read video content.", err)
				return "", err
			}

			// write video file.
			err = os.WriteFile(filePath, videoContent, 0644)
			if err != nil {
				fmt.Println("Can't write video", err)
				return "", err
			}
		case 202:
			fmt.Println("Generation in-progress... automatically try again after 5 sec.")
			time.Sleep(5 * time.Second)
		default:
			//
			fmt.Println(response.Status)
			//
			fmt.Println("Can't connect api.", err)
			var errorMessage map[string]interface{}
			err = json.NewDecoder(response.Body).Decode(&errorMessage)
			if err != nil {
				fmt.Println("Can't request with Non token", err)
				return "", err
			}
		}
	}
	return filePath, nil
}

func deleteFile(path string) {
	if _, err := os.Stat(path); err == nil {
		err := os.Remove(path)
		if err != nil {
			fmt.Println("Fail to delete", err)
		}
	}
}
