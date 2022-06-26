package pop_shark

import (
	"google.golang.org/grpc"
	"log"
	"net"
)

type ChessServer struct {
	port       string
	grpcServer *grpcServer
}

func NewChessServer(port string) *ChessServer {
	return &ChessServer{
		port:       port,
		grpcServer: NewGrpcServer(),
	}
}

func (s ChessServer) Start() {
	lis, err := net.Listen("tcp", s.port)
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
