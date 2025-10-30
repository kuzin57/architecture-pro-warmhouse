package handlers

const (
	// Error codes
	internalServerErrorCode     = "internal_error"
	notFoundErrorCode           = "not_found"
	accessDeniedErrorCode       = "access_denied"
	validationErrorCode         = "validation_error"
	invalidCredentialsErrorCode = "invalid_credentials"
	deviceNotFoundErrorCode     = "device_not_found"

	// Error messages
	internalServerErrorMessage      = "Внутренняя ошибка сервера"
	deviceNotFoundErrorMessage      = "Устройство не найдено или находится в процессе создания"
	houseNotFoundErrorMessage       = "Дом не найден"
	userNotFoundErrorMessage        = "Пользователь не найден"
	accessDeniedErrorMessage        = "Доступ запрещен"
	emailRequiredErrorMessage       = "Email обязателен"
	passwordRequiredErrorMessage    = "Пароль обязателен"
	invalidCredentialsErrorMessage  = "Неверный email или пароль"
	limitValidationErrorMessage     = "limit должен быть от 1 до 1000"
	offsetValidationErrorMessage    = "offset не может быть отрицательным"
	dateRangeValidationErrorMessage = "from_date не может быть больше to_date"
	deviceNotFoundByIDErrorMessage  = "Устройство с указанным ID не найдено"
	deviceAccessDeniedErrorMessage  = "Устройство не принадлежит указанному пользователю"
)
