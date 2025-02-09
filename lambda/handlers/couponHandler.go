package handlers

import (
	"OriD19/webdev2/domain"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func (handler *APIGatewayHandler) GetAllCouponsHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	next := request.QueryStringParameters["next"]

	couponsRange, err := handler.coupons.GetAllCoupons(ctx, &next)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	return Response(200, couponsRange)
}

func (handler *APIGatewayHandler) GetCouponHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	id, ok := request.PathParameters["id"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'id' parameter in path")
	}

	coupon, err := handler.coupons.GetCoupon(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	if coupon == nil {
		return ErrResponse(http.StatusNotFound, "coupon not found")
	}

	return Response(200, coupon)
}

func (handler *APIGatewayHandler) GetAllCouponsFromCategoryHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	category, ok := request.PathParameters["category"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'category' parameter in path")
	}

	couponsRange, err := handler.coupons.GetAllCouponsFromCategory(ctx, category)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	return Response(200, couponsRange)
}

func (handler *APIGatewayHandler) PutCouponHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	id := request.PathParameters["id"]

	if strings.TrimSpace(request.Body) == "" {
		return ErrResponse(http.StatusBadRequest, "missing request body")
	}

	coupon, err := handler.coupons.PutCoupon(ctx, &id, []byte(request.Body))

	if err != nil {
		if errors.Is(err, domain.ErrJsonUnmarshal) {
			return ErrResponse(http.StatusBadRequest, err.Error())
		} else if errors.Is(err, domain.ErrProductIdMismatch) {
			return ErrResponse(http.StatusBadRequest, err.Error())
		}
	}

	return Response(200, coupon)
}

func (handler *APIGatewayHandler) RedeemCouponHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	id, ok := request.PathParameters["id"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'id' parameter in path")
	}

	err := handler.coupons.RedeemCoupon(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	return Response(200, "coupon redeemed successfully")
}

// we can call this handler with a POST request to /coupons/{couponId}/buy/{userId}
func (handler *APIGatewayHandler) BuyCouponHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	couponId, ok := request.PathParameters["id"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'couponId' parameter in path")
	}
	userId, ok := request.PathParameters["userId"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'userId' parameter in path")
	}

	err := handler.coupons.BuyCoupon(ctx, couponId, userId)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	return Response(200, "coupon bought successfully")
}

func (handler *APIGatewayHandler) GetUserOffersHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	id, ok := request.PathParameters["userId"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'id' parameter in path")
	}

	offers, err := handler.coupons.GetUserOffers(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error())
	}

	return Response(200, offers)
}
