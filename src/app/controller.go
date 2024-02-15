package app

import (
	"fmt"
	"net"

	"github.com/robert-min/ArtCore-Go/src/pb"
	"google.golang.org/grpc"
)

func StreamController() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		fmt.Println("Failed to connect: ", err)
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStreamServiceServer(grpcServer, &server{})

	fmt.Println("Server is listening on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Println("Failed to serve: ", err)
	}
}
