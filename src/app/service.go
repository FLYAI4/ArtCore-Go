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

func (s *server) FocusPointStream(req *pb.Request, stream pb.StreamService_FocusPointStreamServer) error {
	// 환경변수
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error to load env.")
	}

	// get env
	openAiToken := os.Getenv("OPEN_AI_KEY")
	// stabilityAiToken := os.Getenv("STABILITY_AI_TOKEN_KEY")

	// TODO : storage 폴더 생성
	storagePath := "./storage"
	userFolderPath := filepath.Join(storagePath, req.Id)

	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		err := os.Mkdir(storagePath, 0755)
		if err != nil {
			fmt.Println("Error to make storage folder. : ", err)
			return err
		}
	}
	// TODO : storage/id 폴더 생성
	if _, err := os.Stat(userFolderPath); os.IsNotExist(err) {
		err := os.Mkdir(userFolderPath, 0755)
		if err != nil {
			fmt.Println("Error to make storage folder. : ", err)
			return err
		}
	}
	// TODO : image 파일 storage/id 폴더 경로에 저장
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
