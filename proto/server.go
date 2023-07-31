package proto

import (
	"context"
	"encoding/json"
	pb "epicpaste/proto/paste"
	"epicpaste/system/model"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedPasteServiceServer
}

func (s *Server) GetPaste(ctx context.Context, in *pb.PasteId) (*pb.Paste, error) {
	log.Printf("%v", in.GetId())
	var paste model.Paste
	var response pb.Paste

	paste.Get(in.GetId())
	data, _ := json.Marshal(paste)
	if err := json.Unmarshal(data, &response); err != nil {
		log.Printf("json unmanshal error%v", err)
	}
	return &response, nil
}

func Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPasteServiceServer(s, &Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
