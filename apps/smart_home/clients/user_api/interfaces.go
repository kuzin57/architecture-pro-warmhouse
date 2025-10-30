package userapi

import (
	"context"
	"smarthome/generated/client"
)

type wrappedClient interface {
	GetDefaultUserWithResponse(ctx context.Context, params *client.GetDefaultUserParams, reqEditors ...client.RequestEditorFn) (*client.GetDefaultUserResponse, error)
}
