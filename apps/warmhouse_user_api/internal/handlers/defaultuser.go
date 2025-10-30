package handlers

import (
	"context"
	"log"

	"github.com/warmhouse/warmhouse_user_api/internal/config"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
)

type GetDefaultUserHandler struct {
	usersService  UsersService
	apiKeyChecker *ApiKeyChecker
}

func NewGetDefaultUserHandler(usersService UsersService, secrets *config.Secrets) *GetDefaultUserHandler {
	return &GetDefaultUserHandler{usersService: usersService, apiKeyChecker: newApiKeyChecker(secrets)}
}

func (h *GetDefaultUserHandler) GetDefaultUser(ctx context.Context, request server.GetDefaultUserRequestObject) (server.GetDefaultUserResponseObject, error) {
	if err := h.apiKeyChecker.CheckApiKey(ctx, request.Params.XApiKey); err != nil {
		return server.GetDefaultUser403JSONResponse{
			Error:   accessDeniedErrorCode,
			Message: accessDeniedErrorMessage,
		}, nil
	}

	resp, err := h.usersService.GetDefaultUser(ctx)
	if err != nil {
		log.Println("Error getting default user:", err)

		return server.GetDefaultUser500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	return server.GetDefaultUser200JSONResponse(resp), nil
}
