package handler

import (
	"context"
	userpb "grpc_and_microservice/proto/user"
	"grpc_and_microservice/user-service/internal/dto"
	"grpc_and_microservice/user-service/internal/service"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {

	input := dto.RegisterDTO{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.service.Register(input)
	if err != nil {
		return nil, err
	}

	return &userpb.RegisterResponse{
		Id:    int64(user.ID),
		Email: user.Email,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {

	input := dto.LoginDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := h.service.Login(input)
	if err != nil {
		return nil, err
	}

	return &userpb.LoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (h *UserHandler) RefreshToken(ctx context.Context, req *userpb.RefreshTokenRequest) (*userpb.RefreshTokenResponse, error) {

	res, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &userpb.RefreshTokenResponse{
		AccessToken: res.AccessToken,
	}, nil
}

func (h *UserHandler) ValidateToken(ctx context.Context, req *userpb.ValidateTokenRequest) (*userpb.ValidateTokenResponse, error) {

	userID, err := h.service.ValidateToken(req.AccessToken)
	if err != nil {
		return &userpb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &userpb.ValidateTokenResponse{
		UserId: int64(userID),
		Valid:  true,
	}, nil
}

func (h *UserHandler) Logout(ctx context.Context, req *userpb.LogoutRequest) (*userpb.LogoutResponse, error) {

	err := h.service.Logout(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &userpb.LogoutResponse{
		Success: true,
	}, nil
}