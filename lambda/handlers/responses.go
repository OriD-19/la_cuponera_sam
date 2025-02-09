package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func Response(code int, bodyObject interface{}) events.APIGatewayProxyResponse {

	// validate the received body

	marshalled, err := json.Marshal(bodyObject)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(marshalled),
		IsBase64Encoded: false,
	}
}

func ErrResponse(code int, body string) events.APIGatewayProxyResponse {
	message := map[string]string{
		"message": body,
	}

	messageBytes, _ := json.Marshal(&message)

	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(messageBytes),
	}
}
