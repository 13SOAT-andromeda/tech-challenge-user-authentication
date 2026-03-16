package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"tech-challenge-user-validation/internal/core/usecases"

	"github.com/aws/aws-lambda-go/events"
)

type AuthHandler struct {
	authUseCase *usecases.AuthUseCase
}

func NewAuthHandler(uc *usecases.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: uc,
	}
}

func (h *AuthHandler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cpf := req.Headers["x-cpf"]
	if cpf == "" {
		return h.errorResponse(http.StatusBadRequest, "x-cpf header is required"), nil
	}

	token, err := h.authUseCase.Authenticate(ctx, cpf)
	if err != nil {
		return h.errorResponse(http.StatusUnauthorized, err.Error()), nil
	}

	body, _ := json.Marshal(map[string]string{
		"token": token,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func (h *AuthHandler) errorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]string{
		"error": message,
	})
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}
