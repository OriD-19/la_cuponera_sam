package types

import (
	"encoding/json"
	"time"
)

type CustomTime struct {
	Time time.Time
}

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
	Title            string     `json:"title" validate:"required"`
	RegularPrice     float32    `json:"regularPrice" validate:"required,gte=0"`
	OfferPrice       float32    `json:"offerPrice" validate:"required,gte=0"`
	AvailableCoupons int        `json:"availableCoupons" validate:"required,gte=1"`
	ExpiresAt        CustomTime `json:"expiresAt" validate:"required"`
	OfferDesc        string     `json:"offerDesc" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validator:"required"`
	Password string `json:"password" validator:"required"`
}

func (c *CustomTime) UnmarshalJSON(data []byte) error {
	// parse the date in a YYYY-MM-DD format
	var timeString string

	if err := json.Unmarshal(data, &timeString); err != nil {
		return err
	}

	newTime, err := time.Parse(DATE_YYYY_MM_DD, timeString)

	if err != nil {
		return err
	}

	c.Time = newTime
	return nil
}
