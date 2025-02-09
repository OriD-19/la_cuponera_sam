package handlers

import (
	"OriD19/webdev2/domain"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func (handler *APIGatewayHandler) GetAllCouponsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	next := request.QueryStringParameters["next"]

	couponsRange, err := handler.coupons.GetAllCoupons(ctx, &next)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	return Response(200, couponsRange), nil
}

func (handler *APIGatewayHandler) GetCouponHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters["id"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'id' parameter in path"), nil
	}

	coupon, err := handler.coupons.GetCoupon(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	if coupon == nil {
		return ErrResponse(http.StatusNotFound, "coupon not found"), nil
	}

	return Response(200, coupon), nil
}

func (handler *APIGatewayHandler) GetAllCouponsFromCategoryHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	category, ok := request.PathParameters["category"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'category' parameter in path"), nil
	}

	couponsRange, err := handler.coupons.GetAllCouponsFromCategory(ctx, category)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	return Response(200, couponsRange), nil
}

func (handler *APIGatewayHandler) PutCouponHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]

	if strings.TrimSpace(request.Body) == "" {
		return ErrResponse(http.StatusBadRequest, "missing request body"), nil
	}

	coupon, err := handler.coupons.PutCoupon(ctx, &id, []byte(request.Body))

	if err != nil {
		if errors.Is(err, domain.ErrJsonUnmarshal) {
			return ErrResponse(http.StatusBadRequest, err.Error()), err
		} else if errors.Is(err, domain.ErrProductIdMismatch) {
			return ErrResponse(http.StatusBadRequest, err.Error()), err
		}
	}

	return Response(200, coupon), nil
}

func (handler *APIGatewayHandler) RedeemCouponHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters["id"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'id' parameter in path"), nil
	}

	err := handler.coupons.RedeemCoupon(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	return Response(200, "coupon redeemed successfully"), nil
}

// we can call this handler with a POST request to /coupons/{couponId}/buy/{userId}
func (handler *APIGatewayHandler) BuyCouponHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	couponId, ok := request.PathParameters["id"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'couponId' parameter in path"), nil
	}
	userId, ok := request.PathParameters["userId"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'userId' parameter in path"), nil
	}

	generatedOffer, err := handler.coupons.BuyCoupon(ctx, couponId, userId)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	return Response(200, generatedOffer), nil
}

func (handler *APIGatewayHandler) GetUserOffersHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters["userId"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'id' parameter in path"), nil
	}

	offers, err := handler.coupons.GetUserOffers(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	return Response(200, offers), nil
}

// retrieve a single offer
func (handler *APIGatewayHandler) GetUserOfferHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	offerId, ok := request.PathParameters["offerId"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'offerId' parameter in path"), nil
	}

	offer, err := handler.coupons.GetGeneratedOffer(ctx, offerId)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	return Response(200, offer), nil
}
