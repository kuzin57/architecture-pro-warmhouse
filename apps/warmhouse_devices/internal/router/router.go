package router

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"slices"

	"github.com/google/uuid"
	"github.com/warmhouse/warmhouse_devices/internal/consts"
	"github.com/warmhouse/warmhouse_devices/internal/entities"
	"github.com/warmhouse/warmhouse_devices/internal/generated/async"
)

var _ async.AppSubscriber = (*Router)(nil)

type Router struct {
	broker *async.AppController

	gatesService   GatesService
	warmingService WarmingService
}

func NewRouter(broker *async.AppController, gatesService GatesService, warmingService WarmingService) *Router {
	return &Router{broker: broker, gatesService: gatesService, warmingService: warmingService}
}

func (r *Router) DeviceCreatedOperationReceived(ctx context.Context, msg async.DeviceCreatedMessageFromDeviceCreatedChannel) error {
	log.Println("got RouteDeviceCreated", msg)

	if err := r.validateDeviceCreatedMessage(&msg); err != nil {
		log.Println("error validating device created message", err)
		return err
	}

	switch msg.Payload.DeviceType {
	case string(consts.DeviceTypeGates):
		if err := r.gatesService.CreateDevice(ctx, msg.Payload); err != nil {
			log.Println("error creating gates device", err)
			return err
		}
	case string(consts.DeviceTypeTemperature):
		if err := r.warmingService.CreateDevice(ctx, msg.Payload); err != nil {
			log.Println("error creating warming device", err)
			return err
		}
	default:
		return nil
	}

	log.Println("device created successfully")

	return nil
}

func (r *Router) DeviceUpdatedOperationReceived(ctx context.Context, msg async.DeviceUpdatedMessageFromDeviceUpdatedChannel) error {
	log.Println("got RouteDeviceUpdated", msg)

	if err := r.validateDeviceUpdatedMessage(&msg); err != nil {
		return err
	}

	switch msg.Payload.Type {
	case string(consts.DeviceTypeGates):
		if err := r.gatesService.UpdateDevice(ctx, msg.Payload); err != nil {
			log.Println("error updating gates device", err)
			return err
		}
	case string(consts.DeviceTypeTemperature):
		if err := r.warmingService.UpdateDevice(ctx, msg.Payload); err != nil {
			log.Println("error updating warming device", err)
			return err
		}
	default:
		return nil
	}

	log.Println("device updated successfully")

	return nil
}

func (r *Router) DeviceDeletedOperationReceived(ctx context.Context, msg async.DeviceDeletedMessageFromDeviceDeletedChannel) error {
	log.Println("got RouteDeviceDeleted", msg)

	deviceID, err := uuid.Parse(msg.Payload.DeviceId)
	if err != nil {
		log.Println("error parsing device id", err)

		return nil
	}

	switch msg.Payload.Type {
	case string(consts.DeviceTypeGates):
		return r.gatesService.DeleteDevice(ctx, deviceID)
	case string(consts.DeviceTypeTemperature):
		return r.warmingService.DeleteDevice(ctx, deviceID)
	default:
		return nil
	}
}

func (r *Router) HouseDeletedOperationReceived(ctx context.Context, msg async.HouseDeletedMessageFromHouseDeletedChannel) error {
	log.Println("got RouteHouseDeleted", msg)

	houseID, err := uuid.Parse(msg.Payload.HouseId)
	if err != nil {
		return fmt.Errorf("invalid house id: %w", err)
	}

	// Удаляем устройства обоих типов
	if err := r.gatesService.DeleteHouseDevices(ctx, houseID); err != nil {
		return err
	}

	return r.warmingService.DeleteHouseDevices(ctx, houseID)
}

func (r *Router) UserDeletedOperationReceived(ctx context.Context, msg async.UserDeletedMessageFromUserDeletedChannel) error {
	log.Println("got RouteUserDeleted", msg)

	userID, err := uuid.Parse(msg.Payload.UserId)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	// Удаляем устройства обоих типов
	if err := r.gatesService.DeleteUserDevices(ctx, userID); err != nil {
		return err
	}

	return r.warmingService.DeleteUserDevices(ctx, userID)
}

func (r *Router) validateDeviceCreatedMessage(msg *async.DeviceCreatedMessageFromDeviceCreatedChannel) error {
	_, err := uuid.Parse(msg.Payload.HouseId)
	if err != nil {
		return fmt.Errorf("invalid house id: %w", err)
	}

	return nil
}

func (r *Router) validateDeviceUpdatedMessage(msg *async.DeviceUpdatedMessageFromDeviceUpdatedChannel) error {
	if msg.Payload.Host != nil {
		_, err := url.Parse(*msg.Payload.Host)
		if err != nil {
			return fmt.Errorf("invalid host: %w", err)
		}
	}

	if msg.Payload.Status != nil && !slices.Contains(entities.GetDeviceStatuses(), *msg.Payload.Status) {
		return fmt.Errorf("invalid status: %s", *msg.Payload.Status)
	}

	return nil
}
