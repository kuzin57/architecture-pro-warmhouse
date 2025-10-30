package handlers

import (
	"context"
	"errors"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
	"github.com/warmhouse/warmhouse_user_api/internal/services/houses"
)

type DeleteHouseHandler struct {
	housesService HousesService
}

func NewDeleteHouseHandler(housesService HousesService) *DeleteHouseHandler {
	return &DeleteHouseHandler{housesService: housesService}
}

func (h *DeleteHouseHandler) DeleteHouse(ctx context.Context, request server.DeleteHouseRequestObject) (server.DeleteHouseResponseObject, error) {
	err := h.housesService.DeleteHouse(ctx, request.Body.HouseId, request.Params.XUserId)
	if err != nil {
		if errors.Is(err, houses.ErrHouseNotFound) {
			return server.DeleteHouse404JSONResponse{
				Error:   notFoundErrorCode,
				Message: houseNotFoundErrorMessage,
			}, nil
		}

		return server.DeleteHouse500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	return server.DeleteHouse204Response{}, nil
}
