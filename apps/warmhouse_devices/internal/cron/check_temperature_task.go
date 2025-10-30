package cron

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/warmhouse/warmhouse_devices/internal/consts"
	"github.com/warmhouse/warmhouse_devices/internal/entities"
)

type CheckTemperatureTask struct {
	devicesRepository   DevicesRepository
	temperatureClient   TemperatureClient
	telemetryRepository TelemetryRepository

	deviceID uuid.UUID
}

func NewCheckTemperatureTask(
	devicesRepository DevicesRepository,
	temperatureClient TemperatureClient,
	telemetryRepository TelemetryRepository,
	deviceID uuid.UUID,
) *CheckTemperatureTask {
	return &CheckTemperatureTask{
		devicesRepository:   devicesRepository,
		temperatureClient:   temperatureClient,
		telemetryRepository: telemetryRepository,
		deviceID:            deviceID,
	}
}

func (t *CheckTemperatureTask) Run(ctx context.Context) error {
	log.Println("running check temperature task for device", t.deviceID)

	device, err := t.devicesRepository.GetDevice(ctx, t.deviceID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	temperature, err := t.temperatureClient.GetTemperature(ctx, device.Host, device.Location)
	if err != nil {
		if err := t.devicesRepository.SetDeviceStatus(ctx, device.ID, consts.DeviceStatusUnknown); err != nil {
			return fmt.Errorf("failed to set device status: %w", err)
		}

		return fmt.Errorf("failed to get temperature: %w", err)
	}

	if err := t.devicesRepository.SetDeviceStatus(ctx, device.ID, consts.DeviceStatusActive); err != nil {
		return fmt.Errorf("failed to set device status: %w", err)
	}

	telemetry := entities.Telemetry{
		DeviceID:  device.ID,
		Timestamp: time.Now(),
		Data: fmt.Sprintf("temperature:%.1f,location:%s,sensor_id:%s",
			temperature.Temperature, temperature.Location, temperature.SensorID),
	}

	if err := t.telemetryRepository.AddTelemetry(ctx, telemetry); err != nil {
		return fmt.Errorf("failed to add telemetry: %w", err)
	}

	log.Println("check temperature task completed")

	return nil
}
