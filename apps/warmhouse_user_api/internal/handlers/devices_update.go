package handlers

import (
	"context"
	"log"

	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
)

type UpdateDeviceHandler struct {
	devicesService DevicesService
}

func NewUpdateDeviceHandler(devicesService DevicesService) *UpdateDeviceHandler {
	return &UpdateDeviceHandler{devicesService: devicesService}
}

func (h *UpdateDeviceHandler) UpdateDevice(ctx context.Context, request server.UpdateDeviceRequestObject) (server.UpdateDeviceResponseObject, error) {
	deviceID, err := h.devicesService.UpdateDevice(ctx, request)
	if err != nil {
		log.Println("error updating device", err)

		return server.UpdateDevice500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	return server.UpdateDevice200JSONResponse(server.UpdateDeviceResponse{
		DeviceId: deviceID,
	}), nil
}
