package main

import (
	"context"
	"io"
	"log"
	"strconv"

	pb "github.com/aleks-papushin/system-monitor/api/gen"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting 'main' func")
	port := 50000
	target := "localhost:" + strconv.Itoa(port)

	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	log.Println("Connection established")
	c := pb.NewStatServiceClient(conn)

	ctx := context.Background()

	req := &pb.StatsRequest{}
	stream, err := c.GetStats(ctx, req)
	if err != nil {
		log.Fatalf("could not get stats: %v", err)
	}

	log.Println("Defined stream")

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("Stream closed by server")
			break
		}
		if err != nil {
			log.Fatalf("error receiving response: %v", err)
		}
		log.Printf("Stats: %v", resp)
	}
}
