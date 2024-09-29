package main

import (
	"context"
	"io"
	"log"
	"os"
	"strconv"

	pb "github.com/aleks-papushin/system-monitor/api/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("Starting 'main' func")

	args := os.Args
	if len(args) < 3 {
		log.Fatalf("Usage: %s <port> <N> <M>", args[0])
	}

	port, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	n, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalf("Invalid N: %v", err)
	}

	m, err := strconv.Atoi(args[3])
	if err != nil {
		log.Fatalf("Invalid M: %v", err)
	}

	target := "localhost:" + strconv.Itoa(port)

	client, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer client.Close()

	log.Println("Connection established")
	c := pb.NewStatServiceClient(client)

	req := &pb.StatsRequest{
		N: int32(n),
		M: int32(m),
	}
	stream, err := c.GetStats(context.Background(), req)
	if err != nil {
		log.Fatalf("could not get stats: %v", err)
	}

	log.Println("Established stream")

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
