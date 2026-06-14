package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/noedaka/go-config-parser/cmd/grpc_server/internal/handler"
	"github.com/noedaka/go-config-parser/cmd/grpc_server/internal/proto"
	"github.com/noedaka/go-config-parser/internal/parser"
	"github.com/noedaka/go-config-parser/internal/service"
	"github.com/noedaka/go-config-parser/internal/service/rules"
	"google.golang.org/grpc"
)

func main() {
	rules := []service.Rule{
		rules.DebugLogRule{},
		rules.PlaintextPasswordRule{},
		rules.ZeroHostRule{},
		rules.TLSDisabledRule{},
		rules.NewWeakAlgorithmRule(),
	}

	parser := parser.YamlJsonParser{}
	handler := handler.NewHandler(rules, parser)
	grpcServer := grpc.NewServer()
	proto.RegisterConfigRecsServiceServer(grpcServer, handler)

	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Starting gRPC server on 9090")

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Received shutdown signal, initiating graceful shutdown...")
	case err := <-serverErr:
		log.Printf("Server error: %v", err)
	}

	// Graceful shutdown
	gracefulCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-gracefulCtx.Done():
		log.Println("Forced shutdown after timeout")
		grpcServer.Stop()
	case <-stopped:
		log.Println("Server stopped gracefully")
	}
}
