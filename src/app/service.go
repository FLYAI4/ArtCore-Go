package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
	"github.com/robert-min/ArtCore-Go/src/focuspoint"
	"github.com/robert-min/ArtCore-Go/src/imagetovideo"
	"github.com/robert-min/ArtCore-Go/src/pb"
)

type server struct {
	pb.StreamServiceServer
}

func (s *server) GeneratedContentStream(req *pb.Request, stream pb.StreamService_GeneratedContentStreamServer) error {
	// get env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error to load env.")
		return nil
	}

	openAiToken := os.Getenv("OPEN_AI_KEY")
	stabilityAiToken := os.Getenv("STABILITY_AI_TOKEN_KEY")

	// make folder
	storagePath := "./storage"
	userFolderPath := filepath.Join(storagePath, req.Id)

	err = makeFolder(storagePath)
	if err != nil {
		fmt.Println("Error to make folder. : ", err)
		return err
	}

	err = makeFolder(userFolderPath)
	if err != nil {
		fmt.Println("Error to make folder. : ", err)
		return err
	}

	// save image to storage/{id}/origin_img.jpg
	err = os.WriteFile(filepath.Join(userFolderPath, "origin_img.jpg"), req.Image, 0644)
	if err != nil {
		fmt.Println("Error to wirte image. : ", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// generate focus point content
	fpm := focuspoint.NewFocusPointManager(userFolderPath, openAiToken)
	go fpm.GenerateFocusPointContent(&wg, stream)

	// generate image to video content
	vm := imagetovideo.NewVideoManager(userFolderPath, stabilityAiToken)
	go vm.GenerateVideoContent(&wg, stream)

	wg.Wait()

	return nil
}

func makeFolder(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
