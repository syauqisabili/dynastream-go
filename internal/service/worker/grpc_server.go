package worker

import (
	"fmt"
	"net"
	"strconv"
	"stream-session-api/internal/conf/network"
	"stream-session-api/internal/service/stream"
	pb "stream-session-api/internal/service/stream/proto"
	"stream-session-api/pkg"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	lis net.Listener
	s   *grpc.Server
)

func InitGrpcServer() error {
	// Get config instance
	conf := network.Get()

	// Grpc server address
	grpcAddr := conf.Grpc.Ip + ":" + strconv.FormatUint(uint64(conf.Grpc.Port), 10)

	var err error
	lis, err = net.Listen("tcp", grpcAddr)
	if err != nil {
		pkg.LogFatal(err.Error())
		return err
	}

	pkg.LogInfo(fmt.Sprintf("gRPC listening on %s...", lis.Addr()))

	opts := []grpc.ServerOption{}
	s = grpc.NewServer(opts...)

	pb.RegisterStreamServiceServer(s, &stream.Server{})
	reflection.Register(s)

	return nil
}

func GrpcServer() {
	if err := s.Serve(lis); err != nil {
		pkg.LogFatal(fmt.Sprintf("Failed to serve: %v", err))

	}
}
