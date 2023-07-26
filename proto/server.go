package proto

import (
	"context"
	pb "epicpaste/proto/paste"
	"epicpaste/system/model"
	"log"
)

type Server struct {
	pb.UnimplementedPasteServiceServer
}

func (s *Server) GetPaste(ctx context.Context, in *pb.PasteId) (*pb.Paste, error) {
	log.Printf("%v", in.GetId())
	var paste model.Paste

	paste.Get(in.GetId())
	return &pb.Paste{
		Id:       paste.ID,
		Content:  paste.Content,
		Public:   paste.Public,
		Language: paste.Languange,
	}, nil
}
