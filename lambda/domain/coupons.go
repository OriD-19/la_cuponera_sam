package domain

import (
	"OriD19/webdev2/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var (
	ErrJsonUnmarshal     = errors.New("failed to parse product from request body")
	ErrProductIdMismatch = errors.New("product ID in path does not match product ID in body")
)

// implementation of the Coupons store for CRUD operations over coupons

type Coupons struct {
	store types.CouponStore
}

func NewCouponsDomain(s types.CouponStore) *Coupons {
	return &Coupons{
		store: s,
	}
}

func (c *Coupons) GetAllCoupons(ctx context.Context, next *string) (types.CouponRange, error) {

	// check if next is just empty spaces
	if next != nil && strings.TrimSpace(*next) == "" {
		next = nil
	}

	couponRange, err := c.store.GetAllCoupons(ctx, next)

	if err != nil {
		return types.CouponRange{}, err
	}

	return couponRange, nil
}

func (c *Coupons) GetAllCouponsFromCategory(ctx context.Context, category string) (types.CouponRange, error) {
	couponRange, err := c.store.GetAllCouponsFromCategory(ctx, category)

	if err != nil {
		return types.CouponRange{}, err
	}

	return couponRange, nil
}

func (c *Coupons) GetCoupon(ctx context.Context, id string) (*types.Coupon, error) {
	coupon, err := c.store.GetCoupon(ctx, id)

	if err != nil {
		return &types.Coupon{}, err
	}

	return &coupon, nil
}

func (c *Coupons) PutCoupon(ctx context.Context, id *string, body []byte) (*types.Coupon, error) {
	couponRequest := types.CreateNewCouponRequest{}

	if err := json.Unmarshal(body, &couponRequest); err != nil {
		return nil, fmt.Errorf("%w", ErrJsonUnmarshal)
	}

	validate := validator.New()
	err := validate.Struct(couponRequest)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	coupon := types.Coupon{}

	// populate the newly created coupon object
	coupon.EntityType = "coupon"
	coupon.Title = couponRequest.Title
	coupon.RegularPrice = couponRequest.RegularPrice
	coupon.OfferPrice = couponRequest.OfferPrice
	coupon.AvailableCoupons = couponRequest.AvailableCoupons
	coupon.ValidFrom = time.Now()
	coupon.ValidUntil = couponRequest.ExpiresAt.Time
	coupon.OfferDesc = couponRequest.OfferDesc

	if id != nil {
		coupon.Id = *id
	} else {
		// assign new UUID to the coupon
		coupon.Id = uuid.New().String()
	}

	err = c.store.PutCoupon(ctx, coupon)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &coupon, nil
}

func (c *Coupons) RedeemCoupon(ctx context.Context, id string) error {
	err := c.store.RedeemCoupon(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

func (c *Coupons) BuyCoupon(ctx context.Context, couponId string, userId string) (*types.GeneratedOffer, error) {
	generatedOffer, err := c.store.BuyCoupon(ctx, couponId, userId)

	if err != nil {
		return nil, err
	}

	return &generatedOffer, nil
}

func (c *Coupons) GetUserOffers(ctx context.Context, id string) (types.OfferRange, error) {
	offerRange, err := c.store.GetUserOffers(ctx, id)

	if err != nil {
		return types.OfferRange{}, err
	}

	return offerRange, nil
}

func (c *Coupons) GetGeneratedOffer(ctx context.Context, id string) (*types.GeneratedOffer, error) {
	offer, err := c.store.GetGeneratedOffer(ctx, id)

	if err != nil {
		return &types.GeneratedOffer{}, err
	}

	return &offer, nil
}
