package handlers

import "OriD19/webdev2/domain"

type APIGatewayHandler struct {
	coupons *domain.Coupons
	users   *domain.Users
}

func NewAPIGatewayHandler(coupons *domain.Coupons, users *domain.Users) *APIGatewayHandler {
	return &APIGatewayHandler{
		coupons: coupons,
		users:   users,
	}
}
