package handlers

import (
	"context"
	"errors"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
	"github.com/warmhouse/warmhouse_user_api/internal/services/users"
)

type GetUserInfoHandler struct {
	usersService UsersService
}

func NewGetUserInfoHandler(usersService UsersService) *GetUserInfoHandler {
	return &GetUserInfoHandler{usersService: usersService}
}

func (h *GetUserInfoHandler) GetUserInfo(ctx context.Context, request server.GetUserInfoRequestObject) (server.GetUserInfoResponseObject, error) {
	user, err := h.usersService.GetUserInfo(ctx, request.Params.XUserId)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return server.GetUserInfo404JSONResponse{
				Error:   notFoundErrorCode,
				Message: userNotFoundErrorMessage,
			}, nil
		}

		return server.GetUserInfo500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	return server.GetUserInfo200JSONResponse(user), nil
}
