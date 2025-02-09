package types

import "context"

/*
	When a coupon is bought, it is associated with a user ID.
	So when we query data again, we know that said coupon is already taken.

	Redeem operations just modify the "redeemed" field inside the database
*/

type CouponStore interface {
	GetAllCoupons(context.Context, *string) (CouponRange, error)
	GetAllCouponsFromCategory(context.Context, string) (CouponRange, error)
	GetCoupon(context.Context, string) (Coupon, error)
	PutCoupon(context.Context, Coupon) error
	RedeemCoupon(context.Context, string) error
	BuyCoupon(context.Context, string, string) error
	GetUserOffers(context.Context, string) (OfferRange, error)
}
