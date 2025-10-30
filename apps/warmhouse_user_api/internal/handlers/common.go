package handlers

import (
	"context"
	"errors"

	"github.com/warmhouse/warmhouse_user_api/internal/config"
)

type ApiKeyChecker struct {
	secrets *config.Secrets
}

func newApiKeyChecker(secrets *config.Secrets) *ApiKeyChecker {
	return &ApiKeyChecker{secrets: secrets}
}

func (c *ApiKeyChecker) CheckApiKey(ctx context.Context, apiKey string) error {
	if apiKey == c.secrets.SmarthomeAPIKey {
		return nil
	}

	return errors.New("invalid api key")
}
