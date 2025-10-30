package users

import "errors"

// Предопределенные ошибки для использования с errors.Is
var (
	// Ошибки пользователей
	ErrUserNotFound       = errors.New("user_not_found")
	ErrUserAlreadyExists  = errors.New("user_already_exists")
	ErrInvalidCredentials = errors.New("invalid_credentials")
)
