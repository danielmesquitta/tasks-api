package service

import (
	"context"

	"github.com/danielmesquitta/tasks-api/internal/app/rpc/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type HealthCheckService struct {
	pb.UnimplementedHealthCheckServiceServer
}

func NewHealthCheckService() *HealthCheckService {
	return &HealthCheckService{}
}

func (s *HealthCheckService) Check(
	context.Context, *emptypb.Empty,
) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
