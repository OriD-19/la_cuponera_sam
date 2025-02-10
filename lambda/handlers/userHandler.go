package handlers

import (
	"OriD19/webdev2/domain"
	"context"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func (handler *APIGatewayHandler) RegisterClient(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// hash the password before registering the user

	client, err := handler.users.RegisterClient(ctx, []byte(request.Body))

	if errors.Is(err, domain.ErrJsonUnmarshal) {
		return ErrResponse(http.StatusBadRequest, "failed to parse client from request body"), err
	} else if err != nil {
		return ErrResponse(http.StatusBadRequest, "username already taken"), err
	}

	return Response(http.StatusOK, client), nil
}

/*
func (handler *APIGatewayHandler) RegisterEnterprise(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	enterprise, err := handler.users.RegisterEnterprise(ctx, []byte(request.Body))

	if errors.Is(err, domain.ErrJsonUnmarshal) {
		return ErrResponse(http.StatusBadRequest, "failed to parse enterprise from request body"), err
	} else if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	return Response(http.StatusOK, enterprise), nil
}

func (handler *APIGatewayHandler) RegisterEmployee(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	employee, err := handler.users.RegisterEmployee(ctx, []byte(request.Body))

	if errors.Is(err, domain.ErrJsonUnmarshal) {
		return ErrResponse(http.StatusBadRequest, "failed to parse employee from request body"), err
	} else if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	return Response(http.StatusOK, employee), nil
}
*/

func (handler *APIGatewayHandler) GetClient(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters["userId"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'id' parameter in path"), nil
	}

	// we use the username as the userId
	client, err := handler.users.GetClient(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	if client == nil {
		return ErrResponse(http.StatusNotFound, "client not found"), nil
	}

	return Response(http.StatusOK, client), nil
}

func (handler *APIGatewayHandler) GetEmployee(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters["id"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'id' parameter in path"), nil
	}

	employee, err := handler.users.GetEmployee(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	if employee == nil {
		return ErrResponse(http.StatusNotFound, "employee not found"), nil
	}

	return Response(http.StatusOK, employee), nil
}

/*
func (handler *APIGatewayHandler) GetEnterprise(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters["id"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'id' parameter in path"), nil
	}

	enterprise, err := handler.users.GetEnterprise(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	if enterprise == nil {
		return ErrResponse(http.StatusNotFound, "enterprise not found"), nil
	}

	return Response(http.StatusOK, enterprise), nil
}

func (handler *APIGatewayHandler) GetAdministrator(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters["id"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'id' parameter in path"), nil
	}

	administrator, err := handler.users.GetAdministrator(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	if administrator == nil {
		return ErrResponse(http.StatusNotFound, "administrator not found"), nil
	}

	return Response(http.StatusOK, administrator), nil
}
*/
