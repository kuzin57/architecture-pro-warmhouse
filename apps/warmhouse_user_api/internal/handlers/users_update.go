package handlers

import (
	"context"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
)

type UpdateUserHandler struct {
	usersService UsersService
}

func NewUpdateUserHandler(usersService UsersService) *UpdateUserHandler {
	return &UpdateUserHandler{usersService: usersService}
}

func (h *UpdateUserHandler) UpdateUser(ctx context.Context, request server.UpdateUserRequestObject) (server.UpdateUserResponseObject, error) {
	user, err := h.usersService.UpdateUser(ctx, request)
	if err != nil {
		return server.UpdateUser500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	return server.UpdateUser200JSONResponse(user), nil
}
