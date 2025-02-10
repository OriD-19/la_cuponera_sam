package main

import (
	"OriD19/webdev2/database"
	"OriD19/webdev2/domain"
	"OriD19/webdev2/handlers"
	"OriD19/webdev2/middleware"
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
		case "/coupons":
			switch request.HTTPMethod {
			case "GET":
				return handler.GetAllCouponsHandler(ctx, request)
			case "POST":
				// TODO add POST method for administrator to upload coupons
				return handler.PutCouponHandler(ctx, request)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: 404,
					Body:       request.Path + " " + request.Resource + ": Not found",
				}, nil
			}
		case "/coupons/category/{category}":
			switch request.HTTPMethod {
			case "GET":
				return handler.GetAllCouponsFromCategoryHandler(ctx, request)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: 404,
					Body:       request.Path + " " + request.Resource + ": Not found",
				}, nil
			}
		case "/coupons/{couponId}":
			switch request.HTTPMethod {
			case "GET":
				return handler.GetCouponHandler(ctx, request)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: 404,
					Body:       request.Path + " " + request.Resource + ": Not found",
				}, nil
			}
		case "/coupons/{couponId}/buy":
			switch request.HTTPMethod {
			case "POST":
				return middleware.ValidateClientJWTMiddleware(ctx, handler.BuyCouponHandler)(ctx, request)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: 404,
					Body:       request.Path + " " + request.Resource + ": Not found",
				}, nil
			}
		case "/offers/allFromUser/{userId}":
			switch request.HTTPMethod {
			case "GET":
				return middleware.ValidateClientJWTMiddleware(ctx, handler.GetUserOffersHandler)(ctx, request)
			// TODO add POST method for administrator to upload offers
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: 404,
					Body:       request.Path + " " + request.Resource + ": Not found",
				}, nil
			}
		case "/offers/{offerId}":
			switch request.HTTPMethod {
			case "GET":
				return middleware.ValidateClientJWTMiddleware(ctx, handler.GetUserOfferHandler)(ctx, request)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: 404,
					Body:       request.Path + " " + request.Resource + ": Not found",
				}, nil
			}
		case "/offers/{offerId}/redeem":
			switch request.HTTPMethod {
			case "POST":
				return middleware.ValidateEmployeeJWTMiddleware(handler.RedeemCouponHandler)(ctx, request)
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
