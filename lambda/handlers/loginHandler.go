package handlers

import (
	"OriD19/webdev2/types"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
)

// for logging out, just delete the token from the client
type LoginRequest struct {
	Email    string `json:"email" validator:"required,email"`
	Password string `json:"password" validator:"required,email"`
}

func (handler *APIGatewayHandler) LoginClient(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)

	if err != nil {
		return ErrResponse(http.StatusBadRequest, "failed to parse credentials from request body"), err
	}

	// validate the login request information
	validate := validator.New()

	err = validate.Struct(loginRequest)

	if err != nil {
		return ErrResponse(http.StatusBadRequest, "invalid email or password"), err
	}

	client, err := handler.users.GetClient(ctx, loginRequest.Email)

	if err != nil {
		return ErrResponse(http.StatusNotFound, err.Error()), err
	}

	if !types.ValidatePassword(client.Password, loginRequest.Password) {
		return ErrResponse(http.StatusUnauthorized, "invalid password"), nil
	}

	token := types.CreateTokenClient(*client)
	successMsg := fmt.Sprintf(`{"access_token": "%s"}`, token)

	// create a new JWT
	return Response(http.StatusOK, successMsg), nil
}

func (handler *APIGatewayHandler) LoginEmployee(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)

	if err != nil {
		return ErrResponse(http.StatusBadRequest, "failed to parse credentials from request body"), err
	}

	// validate the login request information
	validate := validator.New()

	err = validate.Struct(loginRequest)

	if err != nil {
		return ErrResponse(http.StatusBadRequest, "invalid email or password"), err
	}

	employee, err := handler.users.GetEmployee(ctx, loginRequest.Email)

	if err != nil {
		return ErrResponse(http.StatusNotFound, err.Error()), err
	}

	if !types.ValidatePassword(employee.Password, loginRequest.Password) {
		return ErrResponse(http.StatusUnauthorized, "invalid password"), nil
	}

	token := types.CreateTokenEmployee(*employee)
	successMsg := fmt.Sprintf(`{"access_token": "%s"}`, token)

	// create a new JWT
	return Response(http.StatusOK, successMsg), nil

}
