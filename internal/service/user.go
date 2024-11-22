package service

import (
	"context"
	pb "github.com/Hackaton-UDEVS/auth/internal/genproto/auth"
	logger "github.com/Hackaton-UDEVS/auth/internal/logger"
	"github.com/Hackaton-UDEVS/auth/internal/storage/postgres"
	"go.uber.org/zap"
)

type UserService struct {
	storage *postgres.Storage
	pb.UnimplementedAuthServiceServer
}

func NewUserService(db *postgres.Storage) *UserService {
	return &UserService{storage: db}
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	logs, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}
	resp, err := s.storage.Useri.Login(ctx, req)
	if err != nil {
		logs.Error("Error while logging in", zap.Error(err))
		return nil, err
	}
	logs.Info("Successfully logged in")
	return resp, nil
}

func (s *UserService) RegisterUser(ctx context.Context, req *pb.RegisterUserReq) (*pb.RegisterUserRes, error) {
	logs, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}
	resp, err := s.storage.Useri.RegisterUser(ctx, req)
	if err != nil {
		logs.Error("Error while registering user", zap.Error(err))
		return nil, err
	}
	logs.Info("Successfully registered user")
	return resp, nil
}

func (s *UserService) GetUserByID(ctx context.Context, req *pb.GetUserByIDReq) (*pb.GetUserByIDRes, error) {
	logs, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}
	resp, err := s.storage.Useri.GetUserByID(ctx, req)
	if err != nil {
		logs.Error("Error while getting user", zap.Error(err))
		return nil, err
	}
	logs.Info("Successfully retrieved user")
	return resp, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, req *pb.GetAllUserReq) (*pb.GetAllUserRes, error) {
	logs, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}
	resp, err := s.storage.Useri.GetAllUsers(ctx, req)
	if err != nil {
		logs.Error("Error while getting all users", zap.Error(err))
		return nil, err
	}
	logs.Info("Successfully retrieved all users")
	return resp, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.UpdateUserRes, error) {
	logs, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}
	resp, err := s.storage.Useri.UpdateUser(ctx, req)
	if err != nil {
		logs.Error("Error while updating user", zap.Error(err))
		return nil, err
	}
	logs.Info("Successfully updated user")
	return resp, nil
}
