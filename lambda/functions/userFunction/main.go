package main

import (
	"OriD19/webdev2/database"
	"OriD19/webdev2/domain"
	"OriD19/webdev2/handlers"
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	tableName, ok := os.LookupEnv("TABLE_NAME")

	if !ok {
		panic("TABLE_NAME must be set")
	}

	dynamodb := database.NewDynamoDBClient(context.TODO(), tableName)
	couponDomain := domain.NewCouponsDomain(dynamodb)
	usersDomain := domain.NewUsersDomain(dynamodb)
	handler := handlers.NewAPIGatewayHandler(couponDomain, usersDomain)

	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/users/{userId}/profile":
			switch request.HTTPMethod {
			case "GET":
				return handler.GetClient(ctx, request)
				// TODO: Implement PUT/PATCH methods for updating user profile
			}
		case "/users/client/register":
			return handler.RegisterClient(ctx, request)
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	})
}
