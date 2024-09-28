package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	pb "github.com/aleks-papushin/system-monitor/api/gen"
	"github.com/aleks-papushin/system-monitor/internal/collector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedStatServiceServer
}

func (s *server) GetStats(req *pb.StatsRequest, stream pb.StatService_GetStatsServer) error {
	c := collector.GetMacOSStatCollector()
	avgStatChan := c.CollectMacOSStat(int(req.N), int(req.M))
	for stat := range avgStatChan {
		resp := &pb.StatsResponse{
			LoadAverage: stat.LoadAverage,
			UserUsage:   stat.CpuUsage.UserUsage,
			SysUsage:    stat.CpuUsage.SysUsage,
			Idle:        stat.CpuUsage.Idle,
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
	pb.RegisterStatServiceServer(s, &server{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
