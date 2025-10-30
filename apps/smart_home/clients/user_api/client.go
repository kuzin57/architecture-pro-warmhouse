package userapi

import (
	"context"
	"fmt"
	"net/http"
	"smarthome/config"
	"smarthome/generated/client"
	"smarthome/models"

	"github.com/warmhouse/libraries/convert"
)

type Client struct {
	wrappedClient wrappedClient
	secrets       *config.Secrets
}

func NewClient(serverURL string, secrets *config.Secrets) (*Client, error) {
	wrapped, err := client.NewClientWithResponses(serverURL)
	if err != nil {
		return nil, err
	}

	return &Client{wrappedClient: wrapped, secrets: secrets}, nil
}

func (c *Client) GetDefaultUser(ctx context.Context, reqEditors ...client.RequestEditorFn) (models.UserInfo, error) {
	resp, err := c.wrappedClient.GetDefaultUserWithResponse(ctx, &client.GetDefaultUserParams{
		XApiKey: c.secrets.APIKey,
	}, reqEditors...)
	if err != nil {
		return models.UserInfo{}, err
	}

	if resp.StatusCode() != http.StatusOK {
		return models.UserInfo{}, fmt.Errorf("unexpected status code: %d from user-api", resp.StatusCode())
	}

	if resp.JSON200 == nil || resp.JSON200.User == (client.User{}) {
		return models.UserInfo{}, fmt.Errorf("empty response from user-api")
	}

	user := resp.JSON200.User

	return models.UserInfo{
		ID:             user.Id.String(),
		DefaultHouseID: resp.JSON200.DefaultHouseId.String(),
		Name:           user.Name,
		Email:          user.Email,
		Phone:          convert.PointerToValue(user.Phone),
	}, nil
}
