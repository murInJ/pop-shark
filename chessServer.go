package pop_shark

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type ChessServer struct {
	adress     string
	grpcServer *grpcServer
}

func NewChessServer(port string) *ChessServer {
	adress := "0.0.0.0:" + port
	return &ChessServer{
		adress:     adress,
		grpcServer: newGrpcServer(),
	}
}

func (s ChessServer) Start() {
	lis, err := net.Listen("tcp", s.adress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ser := grpc.NewServer()
	RegisterStringServicesServer(ser, s.grpcServer)
	log.Printf("server listening at %v", lis.Addr())
	if err := ser.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
