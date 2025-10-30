package handlers

import (
	"context"
	"errors"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
	"github.com/warmhouse/warmhouse_user_api/internal/services/houses"
)

type GetHouseInfoHandler struct {
	housesService HousesService
}

func NewGetHouseInfoHandler(housesService HousesService) *GetHouseInfoHandler {
	return &GetHouseInfoHandler{housesService: housesService}
}

func (h *GetHouseInfoHandler) GetHouseInfo(ctx context.Context, request server.GetHouseInfoRequestObject) (server.GetHouseInfoResponseObject, error) {
	house, err := h.housesService.GetHouseInfo(ctx, request.Body.HouseId, request.Params.XUserId)
	if err != nil {
		if errors.Is(err, houses.ErrHouseNotFound) {
			return server.GetHouseInfo404JSONResponse{
				Error:   notFoundErrorCode,
				Message: houseNotFoundErrorMessage,
			}, nil
		}

		return server.GetHouseInfo500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	return server.GetHouseInfo200JSONResponse{
		House: house,
	}, nil
}
