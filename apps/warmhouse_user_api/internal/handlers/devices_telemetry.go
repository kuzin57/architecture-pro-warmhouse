package handlers

import (
	"context"
	"errors"
	"time"
	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
	"github.com/warmhouse/warmhouse_user_api/internal/services/devices"
)

type GetDeviceTelemetryHandler struct {
	devicesService DevicesService
}

func NewGetDeviceTelemetryHandler(devicesService DevicesService) *GetDeviceTelemetryHandler {
	return &GetDeviceTelemetryHandler{devicesService: devicesService}
}

func (h *GetDeviceTelemetryHandler) GetDeviceTelemetry(ctx context.Context, request server.GetDeviceTelemetryRequestObject) (server.GetDeviceTelemetryResponseObject, error) {
	// Парсим device_id из пути
	deviceID := request.DeviceId

	// Парсим параметры запроса
	limit := 100 // значение по умолчанию
	if request.Params.Limit != nil {
		limit = *request.Params.Limit
		if limit < 1 || limit > 1000 {
			return server.GetDeviceTelemetry400JSONResponse{
				Error:   validationErrorCode,
				Message: limitValidationErrorMessage,
			}, nil
		}
	}

	offset := 0 // значение по умолчанию
	if request.Params.Offset != nil {
		offset = *request.Params.Offset
		if offset < 0 {
			return server.GetDeviceTelemetry400JSONResponse{
				Error:   validationErrorCode,
				Message: offsetValidationErrorMessage,
			}, nil
		}
	}

	// Парсим даты если они указаны
	var fromDate, toDate *time.Time
	if request.Params.FromDate != nil {
		fromDate = request.Params.FromDate
	}

	if request.Params.ToDate != nil {
		toDate = request.Params.ToDate
	}

	// Проверяем корректность диапазона дат
	if fromDate != nil && toDate != nil && fromDate.After(*toDate) {
		return server.GetDeviceTelemetry400JSONResponse{
			Error:   validationErrorCode,
			Message: dateRangeValidationErrorMessage,
		}, nil
	}

	// Получаем данные телеметрии
	telemetry, err := h.devicesService.GetDeviceTelemetry(ctx, deviceID, request.Params.XUserId, limit, offset, fromDate, toDate)
	if err != nil {
		// Проверяем тип ошибки для возврата соответствующего HTTP статуса
		if errors.Is(err, devices.ErrDeviceNotFound) {
			return server.GetDeviceTelemetry404JSONResponse{
				Error:   deviceNotFoundErrorCode,
				Message: deviceNotFoundByIDErrorMessage,
			}, nil
		}
		if errors.Is(err, devices.ErrDeviceAccessDenied) {
			return server.GetDeviceTelemetry403JSONResponse{
				Error:   accessDeniedErrorCode,
				Message: deviceAccessDeniedErrorMessage,
			}, nil
		}

		// Внутренняя ошибка сервера
		return server.GetDeviceTelemetry500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	return server.GetDeviceTelemetry200JSONResponse(server.TelemetryListResponse{
		Telemetry: telemetry,
	}), nil
}
