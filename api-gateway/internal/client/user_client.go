package client

import (
	"context"

	userpb "grpc_and_microservice/proto/user"
)

type UserClient struct {
	client userpb.UserServiceClient
}

func NewUserClient(client userpb.UserServiceClient) *UserClient {
	return &UserClient{client: client}
}

func (u *UserClient) Register(ctx context.Context, name, email, password string) (*userpb.RegisterResponse, error) {

	return u.client.Register(ctx, &userpb.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})
}

func (u *UserClient) Login(ctx context.Context, email, password string) (*userpb.LoginResponse, error) {

	return u.client.Login(ctx, &userpb.LoginRequest{
		Email:    email,
		Password: password,
	})
}

func (u *UserClient) Refresh(ctx context.Context, token string) (*userpb.RefreshTokenResponse, error) {

	return u.client.RefreshToken(ctx, &userpb.RefreshTokenRequest{
		RefreshToken: token,
	})
}

func (u *UserClient) Validate(ctx context.Context, token string) (*userpb.ValidateTokenResponse, error) {

	return u.client.ValidateToken(ctx, &userpb.ValidateTokenRequest{
		AccessToken: token,
	})
}

func (u *UserClient) Logout(ctx context.Context, token string) (*userpb.LogoutResponse, error) {

	return u.client.Logout(ctx, &userpb.LogoutRequest{
		RefreshToken: token,
	})
}