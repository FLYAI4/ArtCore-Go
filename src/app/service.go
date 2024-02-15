package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/robert-min/ArtCore-Go/src/focuspoint"
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
	// stabilityAiToken := os.Getenv("STABILITY_AI_TOKEN_KEY")

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

	// TODO : focuspoint main content stream
	// TODO : focuspoint coord content stream(각각 함 수 분리?)
	fpm := focuspoint.NewFocusPointManager(userFolderPath, openAiToken)
	response := fpm.GenerateFocusPointContent()
	if err := stream.Send(&pb.Response{Tag: "focus", Data: response}); err != nil {
		fmt.Println("Failed to send response: ", err)
		return err
	}
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
