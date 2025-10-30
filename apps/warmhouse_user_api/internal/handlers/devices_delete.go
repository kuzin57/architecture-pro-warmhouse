package handlers

import (
	"context"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
)

type DeleteDeviceHandler struct {
	devicesService DevicesService
}

func NewDeleteDeviceHandler(devicesService DevicesService) *DeleteDeviceHandler {
	return &DeleteDeviceHandler{devicesService: devicesService}
}

func (h *DeleteDeviceHandler) DeleteDevice(ctx context.Context, request server.DeleteDeviceRequestObject) (server.DeleteDeviceResponseObject, error) {
	err := h.devicesService.DeleteDevice(ctx, request.Body.DeviceId, request.Params.XUserId)
	if err != nil {
		return server.DeleteDevice500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	return server.DeleteDevice204Response{}, nil
}
