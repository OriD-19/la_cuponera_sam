package types

type CouponResponseType struct {
	Coupon
	EnterpriseDetails Enterprise `json:"enterprise"`
}
