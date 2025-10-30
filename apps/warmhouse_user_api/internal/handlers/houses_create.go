package handlers

import (
	"context"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
)

type CreateHouseHandler struct {
	housesService HousesService
}

func NewCreateHouseHandler(housesService HousesService) *CreateHouseHandler {
	return &CreateHouseHandler{housesService: housesService}
}

func (h *CreateHouseHandler) CreateHouse(ctx context.Context, request server.CreateHouseRequestObject) (server.CreateHouseResponseObject, error) {
	houseID, err := h.housesService.CreateHouse(ctx, request)
	if err != nil {
		return server.CreateHouse500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	return server.CreateHouse201JSONResponse{
		Id: houseID,
	}, nil
}
