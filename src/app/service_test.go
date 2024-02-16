package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/robert-min/ArtCore-Go/src/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGeneratedContentStream(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to go server.: %v", err)
	}
	defer conn.Close()

	// given : 유효한 요청
	client := pb.NewStreamServiceClient(conn)
	testImgBytes, err := os.ReadFile("./img/test.jpg")
	if err != nil {
		t.Fatalf("Failed to read test image file. : %v", err)
	}

	request := &pb.Request{Image: testImgBytes, Id: "test1234"}

	// when : focus point 요청
	stream, err := client.GeneratedContentStream(context.Background(), request)
	if err != nil {
		t.Fatalf("Failed to request stream. : %v", err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Failed to get response. : %v", err)
		}
		// then : tag, value 응답
		// assert.True(t, response.Tag == "focus", "Focuspoint 요청에 실패했습니다.")
		fmt.Printf("Tag: %s\n", response.Tag)
		fmt.Printf("Data: %s\n", string(response.Data))
	}

}
