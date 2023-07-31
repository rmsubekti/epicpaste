package main

import (
	"context"
	"log"
	"time"

	pb "epicpaste/proto/paste"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	c := pb.NewPasteServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetPaste(ctx, &pb.PasteId{Id: "f663c9d2-a933-42f2-b57f-efd79dd18fe8"})
	if err != nil {
		log.Fatal("Could not get paste :", err)
	}
	log.Printf("Paste : %s with id %s", r.GetContent(), r.GetId())
}
