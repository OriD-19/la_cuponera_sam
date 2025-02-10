package types

import "time"

type RegisterClientRequest struct {
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	Address     string `json:"address,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	DUI         string `json:"dui" validator:"required"`
}

type RegisterEmployeeRequest struct {
	Username    string `json:"username" validate:"required,email"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	DUI         string `json:"dui" validator:"required"`
}

type CreateNewCouponRequest struct {
	Title            string    `json:"title" validate:"required"`
	RegularPrice     float32   `json:"regularPrice" validate:"required,gte=0"`
	OfferPrice       float32   `json:"offerPrice" validate:"required,gte=0"`
	AvailableCoupons int       `json:"availableCoupons" validate:"required,gte=1"`
	ExpiresAt        time.Time `json:"expiresAt" validate:"required,gt"`
	OfferDesc        string    `json:"offerDesc" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validator:"required"`
	Password string `json:"password" validator:"required"`
}
