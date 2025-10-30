package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/warmhouse/warmhouse_user_api/internal/entities"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"

	"github.com/warmhouse/libraries/convert"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	usersRepository  UsersRepository
	housesRepository HousesRepository
}

func NewService(usersRepository UsersRepository, housesRepository HousesRepository) *Service {
	return &Service{
		usersRepository:  usersRepository,
		housesRepository: housesRepository,
	}
}

func (s *Service) RegisterUser(ctx context.Context, request server.RegisterUserRequestObject) (uuid.UUID, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Body.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error generating password hash: %w", err)
	}

	user := entities.User{
		ID:             uuid.New(),
		Email:          string(request.Body.Email),
		Phone:          convert.ToNullString(request.Body.Phone),
		Name:           request.Body.Name,
		HashedPassword: string(hashedPassword),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error creating user: %w", err)
	}

	return user.ID, nil
}

func (s *Service) UpdateUser(ctx context.Context, request server.UpdateUserRequestObject) (server.User, error) {
	user, err := s.usersRepository.GetUser(ctx, request.Params.XUserId)
	if err != nil {
		return server.User{}, err
	}

	user.Name = convert.UnwrapOr(request.Body.Name, user.Name)
	if request.Body.Phone != nil {
		user.Phone = convert.ToNullString(request.Body.Phone)
	}

	err = s.usersRepository.UpdateUser(ctx, user)
	if err != nil {
		return server.User{}, fmt.Errorf("error updating user: %w", err)
	}

	return server.User{
		Id:        user.ID,
		Email:     user.Email,
		Phone:     convert.FromNullString(user.Phone),
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *Service) GetUserInfo(ctx context.Context, userID uuid.UUID) (server.User, error) {
	user, err := s.usersRepository.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.User{}, ErrUserNotFound
		}
		return server.User{}, err
	}

	return server.User{
		Id:        user.ID,
		Email:     user.Email,
		Phone:     convert.FromNullString(user.Phone),
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *Service) LoginUser(ctx context.Context, request server.LoginUserRequestObject) (server.UserLoginResponse, error) {
	user, err := s.usersRepository.GetUserByEmail(ctx, string(request.Body.Email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.UserLoginResponse{}, ErrInvalidCredentials
		}
		return server.UserLoginResponse{}, fmt.Errorf("error getting user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(request.Body.Password))
	if err != nil {
		return server.UserLoginResponse{}, ErrInvalidCredentials
	}

	return server.UserLoginResponse{
		User: server.User{
			Id:        user.ID,
			Email:     user.Email,
			Phone:     convert.FromNullString(user.Phone),
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (s *Service) GetDefaultUser(ctx context.Context) (server.DefaultUserResponse, error) {
	user, err := s.usersRepository.GetDefaultUser(ctx)
	if err != nil {
		return server.DefaultUserResponse{}, err
	}

	// Получаем самый давний дом пользователя
	house, err := s.housesRepository.GetOldestUserHouse(ctx, user.ID)
	if err != nil {
		return server.DefaultUserResponse{}, fmt.Errorf("error getting default house: %w", err)
	}

	return server.DefaultUserResponse{
		User: server.User{
			Id:        user.ID,
			Email:     user.Email,
			Phone:     convert.FromNullString(user.Phone),
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		DefaultHouseId: house.ID,
	}, nil
}
