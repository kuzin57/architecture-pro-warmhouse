package handlers

import (
	"context"
	"log"

	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
)

type CreateDeviceHandler struct {
	devicesService DevicesService
}

func NewCreateDeviceHandler(devicesService DevicesService) *CreateDeviceHandler {
	return &CreateDeviceHandler{devicesService: devicesService}
}

func (h *CreateDeviceHandler) CreateDevice(ctx context.Context, request server.CreateDeviceRequestObject) (server.CreateDeviceResponseObject, error) {
	deviceID, err := h.devicesService.CreateDevice(ctx, request)
	if err != nil {
		log.Println("error creating device", err)

		return server.CreateDevice500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	return server.CreateDevice201JSONResponse{
		Id: deviceID,
	}, nil
}
