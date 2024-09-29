package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/aleks-papushin/system-monitor/internal/collector"
	"github.com/aleks-papushin/system-monitor/internal/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	gen.UnimplementedStatServiceServer
}

func (s *server) GetStats(req *gen.StatsRequest, stream gen.StatService_GetStatsServer) error {
	c := collector.GetMacOSStatCollector()
	avgStatChan := c.CollectStat(int(req.N), int(req.M))
	for stat := range avgStatChan {
		resp := &gen.StatsResponse{
			LoadAverage: stat.LoadAverage,
			UserUsage:   stat.CPUUsage.UserUsage,
			SysUsage:    stat.CPUUsage.SysUsage,
			Idle:        stat.CPUUsage.Idle,
			Timestamp:   stat.Time.Format(time.RFC3339),
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Usage: %s <port>", args[0])
	}

	port, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	address := ":" + strconv.Itoa(port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	gen.RegisterStatServiceServer(s, &server{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
