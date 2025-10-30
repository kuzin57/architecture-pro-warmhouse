package houses

import "errors"

// Предопределенные ошибки для использования с errors.Is
var (
	// Ошибки домов
	ErrHouseNotFound     = errors.New("house_not_found")
	ErrHouseAccessDenied = errors.New("house_access_denied")
)
