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
		switch request.Resource {
		case "/users/{userId}/profile":
			switch request.HTTPMethod {
			case "GET":
				return handler.GetClient(ctx, request)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: 404,
					Body:       request.Path + " " + request.Resource + ": Not found",
				}, nil
			}
		case "/users/client/register":
			switch request.HTTPMethod {
			case "POST":
				return handler.RegisterClient(ctx, request)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: 404,
					Body:       request.Path + " " + request.Resource + ": Not found",
				}, nil
			}
		default:
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Body:       request.Path + " " + request.Resource + ": Not found",
			}, nil
		}
	})
}
