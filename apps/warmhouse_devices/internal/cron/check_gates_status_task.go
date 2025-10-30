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

type CheckGatesStatusTask struct {
	devicesRepository   DevicesRepository
	gatesClient         GatesClient
	telemetryRepository TelemetryRepository

	deviceID uuid.UUID
}

func NewCheckGatesStatusTask(
	devicesRepository DevicesRepository,
	gatesClient GatesClient,
	telemetryRepository TelemetryRepository,
	deviceID uuid.UUID,
) *CheckGatesStatusTask {
	return &CheckGatesStatusTask{
		devicesRepository:   devicesRepository,
		gatesClient:         gatesClient,
		telemetryRepository: telemetryRepository,
		deviceID:            deviceID,
	}
}

func (t *CheckGatesStatusTask) Run(ctx context.Context) error {
	log.Println("running check gates status task")

	device, err := t.devicesRepository.GetDevice(ctx, t.deviceID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	status, err := t.gatesClient.GetGateStatus(ctx, device.Host)
	if err != nil {
		if err := t.devicesRepository.SetDeviceStatus(ctx, device.ID, consts.DeviceStatusUnknown); err != nil {
			return fmt.Errorf("failed to set device status: %w", err)
		}

		return fmt.Errorf("failed to get gate status: %w", err)
	}

	if status.IsActive {
		if err := t.devicesRepository.SetDeviceStatus(ctx, device.ID, consts.DeviceStatusActive); err != nil {
			return fmt.Errorf("failed to set device status: %w", err)
		}
	} else {
		if err := t.devicesRepository.SetDeviceStatus(ctx, device.ID, consts.DeviceStatusInactive); err != nil {
			return fmt.Errorf("failed to set device status: %w", err)
		}
	}

	telemetry := entities.Telemetry{
		DeviceID:  device.ID,
		Timestamp: time.Now(),
		Data:      status.State.String(),
	}

	if err := t.telemetryRepository.AddTelemetry(ctx, telemetry); err != nil {
		return fmt.Errorf("failed to add telemetry: %w", err)
	}

	log.Println("check gates status task completed")

	return nil
}
