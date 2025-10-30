package gates

import "errors"

var (
	ErrDeviceNotFound    = errors.New("device not found")
	ErrHouseNotFound     = errors.New("house not found")
	ErrAccessDenied      = errors.New("access denied")
	ErrInvalidDeviceData = errors.New("invalid device data")
)
