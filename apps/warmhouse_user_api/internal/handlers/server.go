package handlers

import "github.com/warmhouse/warmhouse_user_api/internal/generated/server"

var (
	_ server.StrictServerInterface = (*Server)(nil)
)

type Server struct {
	*RegisterUserHandler
	*LoginUserHandler
	*GetDefaultUserHandler
	*GetUserInfoHandler
	*UpdateUserHandler
	*GetUserHousesHandler
	*CreateHouseHandler
	*DeleteHouseHandler
	*GetHouseInfoHandler
	*UpdateHouseHandler
	*GetUserDevicesHandler
	*CreateDeviceHandler
	*DeleteDeviceHandler
	*GetDeviceInfoHandler
	*UpdateDeviceHandler
	*GetDeviceTelemetryHandler
}

func NewServer(
	registerUserHandler *RegisterUserHandler,
	loginUserHandler *LoginUserHandler,
	getDefaultUserHandler *GetDefaultUserHandler,
	getUserInfoHandler *GetUserInfoHandler,
	updateUserHandler *UpdateUserHandler,
	getUserHousesHandler *GetUserHousesHandler,
	createHouseHandler *CreateHouseHandler,
	deleteHouseHandler *DeleteHouseHandler,
	getHouseInfoHandler *GetHouseInfoHandler,
	updateHouseHandler *UpdateHouseHandler,
	getUserDevicesHandler *GetUserDevicesHandler,
	createDeviceHandler *CreateDeviceHandler,
	deleteDeviceHandler *DeleteDeviceHandler,
	getDeviceInfoHandler *GetDeviceInfoHandler,
	updateDeviceHandler *UpdateDeviceHandler,
	getDeviceTelemetryHandler *GetDeviceTelemetryHandler,
) *Server {
	return &Server{
		RegisterUserHandler:       registerUserHandler,
		LoginUserHandler:          loginUserHandler,
		GetDefaultUserHandler:     getDefaultUserHandler,
		GetUserInfoHandler:        getUserInfoHandler,
		UpdateUserHandler:         updateUserHandler,
		GetUserHousesHandler:      getUserHousesHandler,
		CreateHouseHandler:        createHouseHandler,
		DeleteHouseHandler:        deleteHouseHandler,
		GetHouseInfoHandler:       getHouseInfoHandler,
		UpdateHouseHandler:        updateHouseHandler,
		GetUserDevicesHandler:     getUserDevicesHandler,
		CreateDeviceHandler:       createDeviceHandler,
		DeleteDeviceHandler:       deleteDeviceHandler,
		GetDeviceInfoHandler:      getDeviceInfoHandler,
		UpdateDeviceHandler:       updateDeviceHandler,
		GetDeviceTelemetryHandler: getDeviceTelemetryHandler,
	}
}
