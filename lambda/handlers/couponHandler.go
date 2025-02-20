package handlers

import (
	"OriD19/webdev2/domain"
	"OriD19/webdev2/types"
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
	id, ok := request.PathParameters["couponId"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'id' parameter in path"), nil
	}

	coupon, err := handler.coupons.GetCoupon(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	var couponResponse types.CouponResponseType

	couponResponse.Coupon = *coupon

	enterprise, err := handler.users.GetEnterprise(ctx, coupon.EnterpriseId)

	if err != nil {
		return ErrResponse(http.StatusNotFound, "enterprise code not found (possibly deleted)"), nil
	}

	couponResponse.EnterpriseDetails = *enterprise

	if coupon == nil {
		return ErrResponse(http.StatusNotFound, "coupon not found"), nil
	}

	return Response(200, couponResponse), nil
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
	var couponId *string

	if strings.TrimSpace(request.Body) == "" {
		return ErrResponse(http.StatusBadRequest, "missing request body"), nil
	}

	if id == "" {
		couponId = nil
	} else {
		couponId = &id
	}

	coupon, err := handler.coupons.PutCoupon(ctx, couponId, []byte(request.Body), handler.users)

	if err != nil {
		if errors.Is(err, domain.ErrJsonUnmarshal) {
			return ErrResponse(http.StatusBadRequest, err.Error()), err
		} else if errors.Is(err, domain.ErrProductIdMismatch) {
			return ErrResponse(http.StatusBadRequest, err.Error()), err
		} else {
			return ErrResponse(http.StatusInternalServerError, err.Error()), err
		}
	}

	return Response(200, coupon), nil
}

func (handler *APIGatewayHandler) RedeemCouponHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := request.PathParameters["offerId"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'offerId' parameter in path"), nil
	}

	// check if the employee is authorized to redeem the coupon

	tokenString := types.ExtractTokenFromHeaders(request.Headers)
	claims, _ := types.ParseToken(tokenString)

	username := claims["username"].(string)
	employee, _ := handler.users.GetEmployee(ctx, username)
	offer, _ := handler.coupons.GetGeneratedOffer(ctx, id)
	coupon, _ := handler.coupons.GetCoupon(ctx, offer.CouponId)

	if employee.EnterpriseId != coupon.EnterpriseId {
		return ErrResponse(http.StatusForbidden, "you must be an employee of this enterprise to redeem this coupon"), nil
	}

	err := handler.coupons.RedeemCoupon(ctx, id)

	if err != nil {
		return ErrResponse(http.StatusBadRequest, err.Error()), nil
	}

	return Response(200, "coupon redeemed successfully"), nil
}

// we can call this handler with a POST request to /coupons/{couponId}/buy/{userId}
func (handler *APIGatewayHandler) BuyCouponHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	couponId, ok := request.PathParameters["couponId"]

	if !ok {
		return ErrResponse(http.StatusBadRequest, "missing 'couponId' parameter in path"), nil
	}

	user, err := types.GetClientAuthFromHeader(request.Headers)

	if err != nil {
		return ErrResponse(http.StatusUnauthorized, err.Error()), err
	}

	// remember: we're using the username as the user id
	generatedOffer, err := handler.coupons.BuyCoupon(ctx, couponId, user.Username)

	if err != nil {
		return ErrResponse(http.StatusInternalServerError, err.Error()), err
	}

	return Response(200, generatedOffer), nil
}

func (handler *APIGatewayHandler) GetUserOffersHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// DONE: extract user id from jwt token
	client, err := types.GetClientAuthFromHeader(request.Headers)

	if err != nil {
		return ErrResponse(401, "client not found"), nil
	}

	offers, err := handler.coupons.GetUserOffers(ctx, client.Username)

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

	// check two cases:
	// - if the user is a client, check that they are the owner of the offer
	// - if the user is an employee, check that they are authorized to view the offer (that is, they are employees of the enterprise issuing the coupon)

	tokenString := types.ExtractTokenFromHeaders(request.Headers)
	claims, _ := types.ParseToken(tokenString)

	userRole := claims["role"].(string)
	username := claims["username"].(string)

	if userRole == "employee" {

		employee, _ := handler.users.GetEmployee(ctx, username)
		coupon, _ := handler.coupons.GetCoupon(ctx, offer.CouponId)

		if coupon.EnterpriseId != employee.EnterpriseId {
			return ErrResponse(http.StatusForbidden, "you must be an employee of this enterprise to access this information"), nil
		}

	} else if userRole == "client" {
		client, _ := handler.users.GetClient(ctx, username)

		if offer.UserId != client.Username {
			return ErrResponse(http.StatusForbidden, "you must be the owner of this offer to view it"), nil
		}

	} else {
		return ErrResponse(http.StatusForbidden, "not authorized for this action"), nil
	}

	return Response(200, offer), nil
}
