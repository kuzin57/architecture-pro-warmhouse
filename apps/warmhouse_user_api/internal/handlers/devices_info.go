package handlers

import (
	"context"
	"errors"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
	"github.com/warmhouse/warmhouse_user_api/internal/services/devices"
)

type GetDeviceInfoHandler struct {
	devicesService DevicesService
}

func NewGetDeviceInfoHandler(devicesService DevicesService) *GetDeviceInfoHandler {
	return &GetDeviceInfoHandler{devicesService: devicesService}
}

func (h *GetDeviceInfoHandler) GetDeviceInfo(ctx context.Context, request server.GetDeviceInfoRequestObject) (server.GetDeviceInfoResponseObject, error) {
	device, err := h.devicesService.GetDeviceInfo(ctx, request.Body.DeviceId, request.Params.XUserId)
	if err != nil {
		if errors.Is(err, devices.ErrDeviceNotFound) {
			return server.GetDeviceInfo404JSONResponse{
				Error:   deviceNotFoundErrorCode,
				Message: deviceNotFoundErrorMessage,
			}, nil
		}

		return server.GetDeviceInfo500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	return server.GetDeviceInfo200JSONResponse{
		Device: device,
	}, nil
}
