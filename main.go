package main

import (
	"fmt"
	"github.com/Hackaton-UDEVS/auth/internal/config"
	"net"

	pb "github.com/Hackaton-UDEVS/auth/internal/genproto/auth"
	logger "github.com/Hackaton-UDEVS/auth/internal/logger"
	"github.com/Hackaton-UDEVS/auth/internal/service"
	"github.com/Hackaton-UDEVS/auth/internal/storage/postgres"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	logs, err := logger.NewLogger()
	if err != nil {
		logs.Error("Error while initializing logger")
		return
	}
	db, err := postgres.ConnectPostgres()
	if err != nil {
		logs.Error("Error while initializing postgres connection")
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.AUTHHOST, cfg.AUTHPORT))
	fmt.Println(cfg.AUTHHOST, cfg.AUTHPORT)
	if err != nil {
		logs.Error("Error while initializing listener")
	}

	defer listener.Close()
	logs.Info(fmt.Sprintf("Server start on port: %d", cfg.AUTHPORT))

	userService := service.NewUserService(db)

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, userService)
	if err := s.Serve(listener); err != nil {
		logs.Error("Error while initializing server")
	}
}
