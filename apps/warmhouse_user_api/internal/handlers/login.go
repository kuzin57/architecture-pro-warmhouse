package handlers

import (
	"context"
	"errors"
	"log"

	"github.com/warmhouse/warmhouse_user_api/internal/generated/server"
	"github.com/warmhouse/warmhouse_user_api/internal/services/users"
)

type LoginUserHandler struct {
	usersService UsersService
}

func NewLoginUserHandler(usersService UsersService) *LoginUserHandler {
	return &LoginUserHandler{usersService: usersService}
}

func (h *LoginUserHandler) LoginUser(ctx context.Context, request server.LoginUserRequestObject) (server.LoginUserResponseObject, error) {
	// Валидация входных данных
	if request.Body.Email == "" {
		return server.LoginUser400JSONResponse{
			Error:   validationErrorCode,
			Message: emailRequiredErrorMessage,
		}, nil
	}

	if request.Body.Password == "" {
		return server.LoginUser400JSONResponse{
			Error:   validationErrorCode,
			Message: passwordRequiredErrorMessage,
		}, nil
	}

	// Вызываем сервис для аутентификации
	response, err := h.usersService.LoginUser(ctx, request)
	if err != nil {
		log.Println("Error logging in user:", err)

		// Проверяем тип ошибки для возврата соответствующего HTTP статуса
		if errors.Is(err, users.ErrInvalidCredentials) {
			return server.LoginUser401JSONResponse{
				Error:   invalidCredentialsErrorCode,
				Message: invalidCredentialsErrorMessage,
			}, nil
		}

		// Внутренняя ошибка сервера
		return server.LoginUser500JSONResponse{
			Error:   internalServerErrorCode,
			Message: internalServerErrorMessage,
		}, nil
	}

	// Возвращаем успешный ответ
	return server.LoginUser200JSONResponse(response), nil
}
