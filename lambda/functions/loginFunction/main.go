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
		case "/login/client":
			switch request.HTTPMethod {
			case "POST":
				return handler.LoginClient(ctx, request)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: 405,
					Body:       "Method not allowed",
				}, nil
			}
		case "/login/employee":
			switch request.HTTPMethod {
			case "POST":
				return handler.LoginEmployee(ctx, request)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: 405,
					Body:       "Method not allowed",
				}, nil
			}
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	})
}
