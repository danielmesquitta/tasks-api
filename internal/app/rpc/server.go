package rpc

import (
	"context"
	"log"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/danielmesquitta/tasks-api/internal/app/rpc/pb"
	"github.com/danielmesquitta/tasks-api/internal/config"
)

func NewServer(
	lc fx.Lifecycle,
	env *config.Env,
	userService pb.UserServiceServer,
	healthService pb.HealthCheckServiceServer,
) *grpc.Server {
	server := grpc.NewServer()

	pb.RegisterUserServiceServer(server, userService)
	pb.RegisterHealthCheckServiceServer(server, healthService)

	reflection.Register(server)

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			lis, err := net.Listen("tcp", ":"+env.Port)
			if err != nil {
				panic(err)
			}

			go func() {
				if err := server.Serve(lis); err != nil {
					panic(err)
				}
			}()

			log.Println("RPC server started on port", env.Port)

			return nil
		},
		OnStop: func(_ context.Context) error {
			server.GracefulStop()
			return nil
		},
	})

	return server
}
