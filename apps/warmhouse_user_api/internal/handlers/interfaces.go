package handlers

import (
	"context"
	"time"

	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"

	"github.com/google/uuid"
)

type DevicesService interface {
	CreateDevice(ctx context.Context, request server.CreateDeviceRequestObject) (uuid.UUID, error)
	GetUserDevices(ctx context.Context, userID uuid.UUID) ([]server.Device, error)
	GetDeviceInfo(ctx context.Context, deviceID, userID uuid.UUID) (server.Device, error)
	UpdateDevice(ctx context.Context, request server.UpdateDeviceRequestObject) (uuid.UUID, error)
	DeleteDevice(ctx context.Context, deviceID, userID uuid.UUID) error
	GetDeviceTelemetry(ctx context.Context, deviceID, userID uuid.UUID, limit, offset int, fromDate, toDate *time.Time) ([]server.Telemetry, error)
}

type HousesService interface {
	GetUserHouses(ctx context.Context, userID uuid.UUID) ([]server.House, error)
	CreateHouse(ctx context.Context, request server.CreateHouseRequestObject) (uuid.UUID, error)
	GetHouseInfo(ctx context.Context, houseID, userID uuid.UUID) (server.House, error)
	UpdateHouse(ctx context.Context, request server.UpdateHouseRequestObject) (server.House, error)
	DeleteHouse(ctx context.Context, houseID, userID uuid.UUID) error
}

type UsersService interface {
	GetUserInfo(ctx context.Context, userID uuid.UUID) (server.User, error)
	UpdateUser(ctx context.Context, request server.UpdateUserRequestObject) (server.User, error)
	RegisterUser(ctx context.Context, request server.RegisterUserRequestObject) (uuid.UUID, error)
	LoginUser(ctx context.Context, request server.LoginUserRequestObject) (server.UserLoginResponse, error)
	GetDefaultUser(ctx context.Context) (server.DefaultUserResponse, error)
}
