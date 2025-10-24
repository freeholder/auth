package main

import (
	"context"
	"log"
	"time"

	"github.com/fatih/color"
	desc "github.com/freeholder/auth/pkg/note_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "172.26.112.94:50051"
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: #{err}")
	}
	defer conn.Close()

	c := desc.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: 1.0})
	if err != nil {
		log.Fatalf("failed to get note by id")
	}

	log.Printf(color.RedString("Note: %v"), r)

}
