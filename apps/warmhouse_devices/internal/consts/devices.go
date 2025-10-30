package consts

type DeviceType string

const (
	DeviceTypeGates       DeviceType = "gates"
	DeviceTypeTemperature DeviceType = "temperature"
)

type DeviceStatus string

const (
	DeviceStatusInactive DeviceStatus = "inactive"
	DeviceStatusActive   DeviceStatus = "active"
	DeviceStatusUnknown  DeviceStatus = "unknown"
)

const (
	DeviceUnitGates   = "-"
	DeviceUnitWarming = "°C"
)

const (
	GatesValueOpened = "opened"
	GatesValueClosed = "closed"
)
