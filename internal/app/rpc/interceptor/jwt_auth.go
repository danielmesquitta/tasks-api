package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type key string

const ClaimsKey key = "claims"

func (i *Interceptor) UnaryEnsureJWTAuthentication(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	log.Println("--> unary interceptor: ", info.FullMethod)

	if err := i.ensureJWTAuthentication(&ctx, info.FullMethod); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func (i *Interceptor) StreamEnsureJWTAuthentication(
	srv any,
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Println("--> stream interceptor: ", info.FullMethod)

	ctx := stream.Context()
	if err := i.ensureJWTAuthentication(&ctx, info.FullMethod); err != nil {
		return err
	}

	return handler(srv, stream)
}

func (i *Interceptor) ensureJWTAuthentication(
	ctx *context.Context,
	method string,
) error {
	allowedRoles, ok := i.AllowedRolesByMethod[method]
	if everyoneCanAccess := !ok; everyoneCanAccess {
		return nil
	}

	md, ok := metadata.FromIncomingContext(*ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(
			codes.Unauthenticated,
			"authorization token is not provided",
		)
	}

	accessToken := values[0]

	claims, err := i.jwt.ValidateAccessToken(accessToken)
	if err != nil {
		return status.Errorf(
			codes.Unauthenticated,
			"invalid or expired token",
		)
	}

	isRoleAllowed := false
	for _, role := range allowedRoles {
		if role == claims.Role {
			isRoleAllowed = true
			break
		}
	}

	if !isRoleAllowed {
		return status.Errorf(
			codes.PermissionDenied,
			"insufficient permissions",
		)
	}

	*ctx = context.WithValue(*ctx, ClaimsKey, claims)

	return nil
}
