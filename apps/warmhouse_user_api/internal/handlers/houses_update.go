package handlers

import (
	"context"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
)

type UpdateHouseHandler struct {
	housesService HousesService
}

func NewUpdateHouseHandler(housesService HousesService) *UpdateHouseHandler {
	return &UpdateHouseHandler{housesService: housesService}
}

func (h *UpdateHouseHandler) UpdateHouse(ctx context.Context, request server.UpdateHouseRequestObject) (server.UpdateHouseResponseObject, error) {
	house, err := h.housesService.UpdateHouse(ctx, request)
	if err != nil {
		return server.UpdateHouse500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	return server.UpdateHouse200JSONResponse(house), nil
}
