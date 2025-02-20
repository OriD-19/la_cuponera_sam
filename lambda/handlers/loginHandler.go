package handlers

import (
	"OriD19/webdev2/types"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
)

// for logging out, just delete the token from the client

func (handler *APIGatewayHandler) LoginClient(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var loginRequest types.LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)

	if err != nil {
		return ErrResponse(http.StatusBadRequest, "failed to parse credentials from request body"), nil
	}

	// validate the login request information
	validate := validator.New()

	err = validate.Struct(loginRequest)

	if err != nil {
		return ErrResponse(http.StatusBadRequest, "invalid username or password"), nil
	}

	client, err := handler.users.GetClient(ctx, loginRequest.Username)

	if err != nil {
		return ErrResponse(http.StatusNotFound, err.Error()), nil
	}

	if !types.ValidatePassword(client.Password, loginRequest.Password) {
		return ErrResponse(http.StatusUnauthorized, "invalid password"), nil
	}

	token := types.CreateTokenClient(*client)

	type LoginClientResponse struct {
		AuthToken string       `json:"authToken"`
		Client    types.Client `json:"client"`
	}

	lcRes := LoginClientResponse{
		AuthToken: token,
		Client:    *client,
	}

	// create a new JWT
	return Response(http.StatusOK, lcRes), nil
}

func (handler *APIGatewayHandler) LoginEmployee(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var loginRequest types.LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)

	if err != nil {
		return ErrResponse(http.StatusBadRequest, "failed to parse credentials from request body"), nil
	}

	// validate the login request information
	validate := validator.New()

	err = validate.Struct(loginRequest)

	if err != nil {
		return ErrResponse(http.StatusBadRequest, "invalid username or password"), nil
	}

	employee, err := handler.users.GetEmployee(ctx, loginRequest.Username)

	if err != nil {
		return ErrResponse(http.StatusNotFound, err.Error()), nil
	}

	if !types.ValidatePassword(employee.Password, loginRequest.Password) {
		return ErrResponse(http.StatusUnauthorized, "invalid password"), nil
	}

	token := types.CreateTokenEmployee(*employee)

	type LoginEmployeeReponse struct {
		AuthToken string         `json:"authToken"`
		Employee  types.Employee `json:"employee"`
	}

	leRes := LoginEmployeeReponse{
		AuthToken: token,
		Employee:  *employee,
	}

	// create a new JWT
	return Response(http.StatusOK, leRes), nil
}
