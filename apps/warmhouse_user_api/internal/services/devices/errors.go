package devices

import "errors"

// Предопределенные ошибки для использования с errors.Is
var (
	// Ошибки устройств
	ErrDeviceNotFound     = errors.New("device_not_found")
	ErrDeviceAccessDenied = errors.New("device_access_denied")
	ErrHouseNotFound      = errors.New("house_not_found")
	ErrHouseAccessDenied  = errors.New("house_access_denied")
)
