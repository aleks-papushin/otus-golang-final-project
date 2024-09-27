package main

import (
	"log"
	"net"
	"strconv"
	"time"

	pb "github.com/aleks-papushin/system-monitor/api/gen"
	"github.com/aleks-papushin/system-monitor/internal/collector"
	"github.com/aleks-papushin/system-monitor/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedStatServiceServer
	statChan <-chan *models.Stat
}

func (s *server) GetStats(req *pb.StatsRequest, stream pb.StatService_GetStatsServer) error {
	for stat := range s.statChan {
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
	port := 50000
	address := ":" + strconv.Itoa(port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	statChan := make(chan *models.Stat, 100)
	pb.RegisterStatServiceServer(s, &server{statChan: statChan})
	reflection.Register(s)

	go func() {
		collector.CollectMacOSStat(statChan)
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
