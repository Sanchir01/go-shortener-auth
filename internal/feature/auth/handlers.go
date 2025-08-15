package auth

import (
	"context"
	authv1 "github.com/Sanchir01/go-shortener-proto/pkg/gen/go/v1/auth"
	"google.golang.org/grpc"
)

type Handlers struct {
	authv1.UnimplementedAuthServer
}

func NewServer(grpcServer *grpc.Server) {
	authv1.RegisterAuthServer(grpcServer, &Handlers{})
}
func (h *Handlers) Register(ctx context.Context, request *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) ConfirmRegister(ctx context.Context, request *authv1.ConfirmRegisterRequest) (*authv1.ConfirmRegisterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	panic("implement me")
}
