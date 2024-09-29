package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/aleks-papushin/system-monitor/internal/proto/gen"
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

	n, err := strconv.ParseInt(args[2], 10, 32)
	if err != nil {
		log.Fatalf("Invalid N: %v", err)
	}

	m, err := strconv.ParseInt(args[3], 10, 32)
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
	c := gen.NewStatServiceClient(client)

	req := &gen.StatsRequest{
		N: n,
		M: m,
	}
	stream, err := c.GetStats(context.Background(), req)
	if err != nil {
		log.Printf("could not get stats: %v", err)
		client.Close()
		return
	}

	log.Println("Established stream")

	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Println("Stream closed by server")
			break
		}
		if err != nil {
			log.Printf("error receiving response: %v", err)
			client.Close()
			return
		}
		log.Printf("Stats: %v", resp)
	}
}
