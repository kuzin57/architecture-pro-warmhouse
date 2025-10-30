package handlers

import (
	"context"

	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
)

type RegisterUserHandler struct {
	usersService UsersService
}

func NewRegisterUserHandler(usersService UsersService) *RegisterUserHandler {
	return &RegisterUserHandler{usersService: usersService}
}

func (h *RegisterUserHandler) RegisterUser(ctx context.Context, request server.RegisterUserRequestObject) (server.RegisterUserResponseObject, error) {
	userID, err := h.usersService.RegisterUser(ctx, request)
	if err != nil {
		return server.RegisterUser500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	return server.RegisterUser201JSONResponse{
		UserId: userID,
	}, nil
}
