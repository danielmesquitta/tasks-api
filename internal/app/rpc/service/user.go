package service

import (
	"context"

	"github.com/danielmesquitta/tasks-api/internal/app/rpc/pb"
	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/domain/usecase"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	createUserUseCase *usecase.CreateUser
}

func NewUserService(
	createUserUseCase *usecase.CreateUser,
) *UserService {
	return &UserService{
		createUserUseCase: createUserUseCase,
	}
}

func (s *UserService) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*emptypb.Empty, error) {
	params := usecase.CreateUserParams{}

	if err := copier.Copy(&params, req); err != nil {
		return nil, entity.NewErr(err)
	}

	if err := s.createUserUseCase.Execute(ctx, params); err != nil {
		return nil, entity.NewErr(err)
	}

	return &emptypb.Empty{}, nil
}
