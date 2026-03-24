package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"tech-challenge-user-validation/internal/core/ports"
	"tech-challenge-user-validation/internal/core/usecases"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
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
	s, _ := tracer.StartSpanFromContext(ctx, "child.span")
	defer s.Finish()
	switch {
	case req.HTTPMethod == http.MethodPost && req.Path == "/sessions":
		return h.handleLogin(ctx, req)
	case req.HTTPMethod == http.MethodPost && req.Path == "/sessions/refresh":
		return h.handleRefresh(ctx, req)
	case req.HTTPMethod == http.MethodDelete && req.Path == "/sessions/logout":
		return h.handleLogout(ctx, req)
	default:
		return h.errorResponse(http.StatusNotFound, "route not found"), nil
	}
}

func (h *AuthHandler) handleLogin(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input ports.LoginInput

	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return h.errorResponse(http.StatusBadRequest, "invalid request body"), nil
	}

	if input.Document == "" {
		return h.errorResponse(http.StatusBadRequest, "document is required"), nil
	}
	if input.Password == "" {
		return h.errorResponse(http.StatusBadRequest, "password is required"), nil
	}

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
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func (h *AuthHandler) handleRefresh(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input ports.RefreshInput

	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return h.errorResponse(http.StatusBadRequest, "invalid request body"), nil
	}
	if input.RefreshToken == "" {
		return h.errorResponse(http.StatusBadRequest, "refresh_token is required"), nil
	}

	output, err := h.authUseCase.Refresh(ctx, input)
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
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func (h *AuthHandler) handleLogout(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	authHeader := req.Headers["Authorization"]
	if authHeader == "" {
		authHeader = req.Headers["authorization"]
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) {
		return h.errorResponse(http.StatusUnauthorized, "missing or invalid Authorization header"), nil
	}
	tokenString := authHeader[len(bearerPrefix):]

	if err := h.authUseCase.Logout(ctx, tokenString); err != nil {
		return h.errorResponse(http.StatusUnauthorized, err.Error()), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
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
