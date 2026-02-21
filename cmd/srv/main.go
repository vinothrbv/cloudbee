package main

import (
	"crud/app/domain/repository"
	"crud/app/external/server"
	"crud/pb"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	repo := repository.NewInMemoryPostRepository()
	srv := server.NewServer(repo)

	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPostServiceServer(grpcServer, srv)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("server stopped gracefully")
}
