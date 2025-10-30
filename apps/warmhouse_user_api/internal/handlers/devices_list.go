package handlers

import (
	"context"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
)

type GetUserDevicesHandler struct {
	devicesService DevicesService
}

func NewGetUserDevicesHandler(devicesService DevicesService) *GetUserDevicesHandler {
	return &GetUserDevicesHandler{devicesService: devicesService}
}

func (h *GetUserDevicesHandler) GetUserDevices(ctx context.Context, request server.GetUserDevicesRequestObject) (server.GetUserDevicesResponseObject, error) {
	devices, err := h.devicesService.GetUserDevices(ctx, request.Params.XUserId)
	if err != nil {
		return server.GetUserDevices500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	return server.GetUserDevices200JSONResponse{
		Devices: devices,
	}, nil
}
