package handlers

import (
	"context"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
)

type GetUserHousesHandler struct {
	housesService HousesService
}

func NewGetUserHousesHandler(housesService HousesService) *GetUserHousesHandler {
	return &GetUserHousesHandler{housesService: housesService}
}

func (h *GetUserHousesHandler) GetUserHouses(ctx context.Context, request server.GetUserHousesRequestObject) (server.GetUserHousesResponseObject, error) {
	houses, err := h.housesService.GetUserHouses(ctx, request.Params.XUserId)
	if err != nil {
		return server.GetUserHouses500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	return server.GetUserHouses200JSONResponse{
		Houses: houses,
	}, nil
}
