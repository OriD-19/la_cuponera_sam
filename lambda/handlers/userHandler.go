package handlers

import (
	"OriD19/webdev2/domain"
	"context"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func (handler *APIGatewayHandler) RegisterClient(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	// hash the password before registering the user

	client, err := handler.users.RegisterClient(ctx, []byte(request.Body))

	if errors.Is(err, domain.ErrJsonUnmarshal) {
		return ErrResponse(http.StatusBadRequest, "failed to parse client from request body")
	} else if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	return Response(http.StatusOK, client)
}

func (handler *APIGatewayHandler) RegisterEnterprise(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {

	enterprise, err := handler.users.RegisterEnterprise(ctx, []byte(request.Body))

	if errors.Is(err, domain.ErrJsonUnmarshal) {
		return ErrResponse(http.StatusBadRequest, "failed to parse enterprise from request body")
	} else if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	return Response(http.StatusOK, enterprise)
}

func (handler *APIGatewayHandler) RegisterEmployee(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {

	employee, err := handler.users.RegisterEmployee(ctx, []byte(request.Body))

	if errors.Is(err, domain.ErrJsonUnmarshal) {
		return ErrResponse(http.StatusBadRequest, "failed to parse employee from request body")
	} else if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	return Response(http.StatusOK, employee)
}

func (handler *APIGatewayHandler) GetClient(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	id, ok := request.PathParameters["id"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'id' parameter in path")
	}

	client, err := handler.users.GetClient(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	if client == nil {
		return ErrResponse(http.StatusNotFound, "client not found")
	}

	return Response(http.StatusOK, client)
}

func (handler *APIGatewayHandler) GetEnterprise(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	id, ok := request.PathParameters["id"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'id' parameter in path")
	}

	enterprise, err := handler.users.GetEnterprise(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	if enterprise == nil {
		return ErrResponse(http.StatusNotFound, "enterprise not found")
	}

	return Response(http.StatusOK, enterprise)
}

func (handler *APIGatewayHandler) GetEmployee(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	id, ok := request.PathParameters["id"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'id' parameter in path")
	}

	employee, err := handler.users.GetEmployee(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	if employee == nil {
		return ErrResponse(http.StatusNotFound, "employee not found")
	}

	return Response(http.StatusOK, employee)
}

func (handler *APIGatewayHandler) GetAdministrator(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	id, ok := request.PathParameters["id"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'id' parameter in path")
	}

	administrator, err := handler.users.GetAdministrator(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	if administrator == nil {
		return ErrResponse(http.StatusNotFound, "administrator not found")
	}

	return Response(http.StatusOK, administrator)
}
