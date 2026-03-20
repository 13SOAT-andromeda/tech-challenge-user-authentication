package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"tech-challenge-user-validation/internal/core/ports"
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
	var input ports.LoginInput

	// 1) Parse body JSON
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return h.errorResponse(http.StatusBadRequest, "invalid request body"), nil
	}

	// 2) Validate required fields
	if input.Document == "" {
		return h.errorResponse(http.StatusBadRequest, "document is required"), nil
	}
	if input.Password == "" {
		return h.errorResponse(http.StatusBadRequest, "password is required"), nil
	}

	// 3) Execute login use case
	output, err := h.authUseCase.Login(ctx, input)
	if err != nil {
		return h.errorResponse(http.StatusUnauthorized, err.Error()), nil
	}

	body, err := json.Marshal(output)
	if err != nil {
		return h.errorResponse(http.StatusInternalServerError, "failed to serialize response"), nil
	}

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
