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
		case "/coupons":
			switch request.HTTPMethod {
			case "GET":
				return handler.GetAllCouponsHandler(ctx, request)
				// TODO add POST method for administrator to upload coupons
			}
		case "/coupons/category/{category}":
			return handler.GetAllCouponsFromCategoryHandler(ctx, request)
		case "/coupons/{id}":
			return handler.GetCouponHandler(ctx, request)
		case "/coupons/buy":
			return handler.BuyCouponHandler(ctx, request)
		case "/offers/{offerId}/redeem":
			return handler.RedeemCouponHandler(ctx, request)
		default:
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Body:       "Not found",
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	})
}
