package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	server "github.com/1055373165/Distributed_KV_Store/distributekv"
	pb "github.com/1055373165/Distributed_KV_Store/grpc/groupcachepb"
	"google.golang.org/grpc"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8001, "port")
	flag.Parse()
	addr := fmt.Sprintf("localhost:%d", port)

	//关闭信号处理
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		s := <-ch
		etcdUnRegister(addr)
		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()

	err := etcdRegister(addr)
	if err != nil {
		panic(err)

	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptor()))
	pb.RegisterGroupCacheServer(grpcServer, &server.Server{})
	log.Printf("service start port %d\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Printf("call %s\n", info.FullMethod)
		resp, err = handler(ctx, req)
		return resp, err
	}
}
