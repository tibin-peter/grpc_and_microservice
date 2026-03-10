package service

import (
	"context"
	"errors"
	"grpc_and_microservice/user-service/internal/dto"
	"grpc_and_microservice/user-service/internal/model"
	"grpc_and_microservice/user-service/internal/repository"
	"grpc_and_microservice/user-service/internal/utilities"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserService interface {
	Register(input dto.RegisterDTO) (*model.User, error)
	Login(input dto.LoginDTO) (*dto.AuthResponseDTO, error)
	RefreshToken(refreshToken string) (*dto.AuthResponseDTO, error)
	ValidateToken(accessToken string) (uint, error)
	Logout(refreshToken string) error
}

type userService struct {
	repo  repository.Repository
	redis *redis.Client
}

func NewUserService(repo repository.Repository, redis *redis.Client) UserService {
	return &userService{
		repo:  repo,
		redis: redis,
	}
}

func (s *userService) Register(input dto.RegisterDTO) (*model.User, error) {

	hash, err := utilities.GenereatePassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hash,
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(input dto.LoginDTO) (*dto.AuthResponseDTO, error) {

	user, err := s.repo.FindUserByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = utilities.CheckPassword(user.Password,input.Password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	jwtKey := os.Getenv("JWT_SECRET")

	accessToken, _, err := utilities.GenerateAccessToken(user.ID, user.Email, "user", jwtKey)
	if err != nil {
		return nil, err
	}

	refreshToken, _, err := utilities.GenerateRefreshToken(user.ID, user.Email, "user", jwtKey)
	if err != nil {
		return nil, err
	}

	err = s.redis.Set(context.Background(), refreshToken, user.ID, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) RefreshToken(refreshToken string) (*dto.AuthResponseDTO, error) {

	// userID, err := s.redis.Get(context.Background(), refreshToken).Result()
	// if err != nil {
	// 	return nil, errors.New("invalid refresh token")
	// }

	jwtKey := os.Getenv("JWT_SECRET")

	claims, err := utilities.ValidateToken(refreshToken, jwtKey)
	if err != nil {
		return nil, err
	}

	newAccessToken, _, err := utilities.GenerateAccessToken(claims.UserID, claims.Email, claims.Role, jwtKey)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponseDTO{
		AccessToken: newAccessToken,
	}, nil
}

func (s *userService) ValidateToken(accessToken string) (uint, error) {

	jwtKey := os.Getenv("JWT_SECRET")

	claims, err := utilities.ValidateToken(accessToken, jwtKey)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}

func (s *userService) Logout(refreshToken string) error {
	return s.redis.Del(context.Background(), refreshToken).Err()
}