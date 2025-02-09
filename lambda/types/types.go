package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const SECRET = "secret"

// ************************************************************
// USER ENTITIES
// ************************************************************
type Entity struct {
	EntityType string `dynamodbav:"entityType" json:"entityType" validator:"oneof=client administrator enterprise employee offer"`
}

type User struct {
	Entity
	Email     string    `dynamodbav:"id" json:"email" validator:"required,email"` // email is the id
	Username  string    `dynamodbav:"username" json:"username" validator:"required,max=100"`
	Password  string    `dynamodbav:"password" json:"password" validator:"required,min=8"`
	CreatedAt time.Time `dynamodbav:"createdAt" json:"createdAt" validator:"required"`

	// a single user can have many coupons, and a single coupon can be owned by many users
	// also, to keep a history of the coupons that a user has redeemed, we won't delete the registers. Just update the original item
	//UserCoupons       []UserCoupon       `dynamodbav:"userCoupons" json:"userCoupons" validator:"required_with=ClientDetails"` // for storing many-to-many relationship
	//RegisteredCoupons []RegisteredCoupon `dynamodbav:"registeredCoupons" json:"registeredCoupons" validator:"omitempty,required_with=EnterpriseDetails"`
}

type Client struct {
	User
	FirstName   string `dynamodbav:"firstName" json:"firstName"`
	LastName    string `dynamodbav:"lastName" json:"lastName"`
	Address     string `dynamodbav:"address" json:"address"`
	PhoneNumber string `dynamodbav:"phoneNumber" json:"phoneNumber"`
	DUI         string `dynamodbav:"dui" json:"dui"`
}

type Enterprise struct {
	User
	EnterpriseCode      string `dynamodbav:"enterpriseCode" json:"enterpriseCode"` // code that is generated in some type of way...
	EnterpriseName      string `dynamodbav:"enterpriseName" json:"enterpriseName"`
	ScheduleDescription string `dynamodbav:"scheduleDescription" json:"scheduleDescription"`
	Location            string `dynamodbav:"location" json:"location"`
	PhoneNumber         string `dynamodbav:"phoneNumber" json:"phoneNumber"`
	Category            string `dynamodbav:"category" json:"category"` // the category of the enterprise (restaurant, gym, etc)
}

type Administrator struct {
	User
	FirstName string `dynamodbav:"firstName" json:"firstName"`
	LastName  string `dynamodbav:"lastName" json:"lastName"`
}

type Employee struct {
	User
	FirstName   string `dynamodbav:"firstName" json:"firstName"`
	LastName    string `dynamodbav:"lastName" json:"lastName"`
	PhoneNumber string `dynamodbav:"phoneNumber" json:"phoneNumber"`
	DUI         string `dynamodbav:"dui" json:"dui"`
}

// ************************************************************
// COUPON ENTITIES
// ************************************************************

type RegisteredCoupon struct {
	CouponId string `dynamodbav:"couponId" json:"couponId"`
	Status   string `dynamodbav:"status" json:"status"`
}

type CouponDetails struct {
	Location            string `dynamodbav:"location" json:"location"`
	ScheduleDescription string `dynamodbav:"scheduleDescription" json:"scheduleDescription"`
	PhoneNumber         string `dynamodbav:"phoneNumber" json:"phoneNumber"`
}

type Coupon struct {
	Entity
	Id           string    `dynamodbav:"id" json:"id"`
	Title        string    `dynamodbav:"title" json:"title" validate:"required,max=100"`
	RegularPrice float32   `dynamodbav:"regularPrice" json:"regularPrice" validate:"required"`
	OfferPrice   float32   `dynamodbav:"offerPrice" json:"offerPrice" validate:"required"`
	ValidFrom    time.Time `dynamodbav:"validFrom" json:"validFrom" validate:"required"`
	ValidUntil   time.Time `dynamodbav:"validUntil" json:"validUntil" validate:"required,gt"` // greater than now

	// Available quantity of coupons. -1 if there is no limit in the amount of coupons
	AvailableCoupons int    `dynamodbav:"availableCoupons" json:"availableCoupons" validate:"required,ne=0"`
	OfferDesc        string `dynamodbav:"offerDesc" json:"offerDesc"`

	// TODO: could use this one, but not sure yet
	// CouponState      string `dynamodbav:"couponState" json:"couponState" validate:"required,oneof=active inactive expired pending rejected"`

	// Can have any structure for other properties defined in the UI
	// the idea is that is a nested details object
	Details      CouponDetails `dynamodbav:"details" json:"details,omitempty" validate:"omitempty"`
	EnterpriseId string        `dynamodbav:"enterpriseCode" json:"enterpriseCode" validate:"required"`
	Category     string        `dynamodbav:"category" json:"category" validate:"required"` // same as the enterprise category defined above
}

// struct for a generated coupon (the one that user buys, not the general offer)
type GeneratedOffer struct {
	Entity
	Id             string    `dynamodbav:"id" json:"id"` // this id will be the generated token for the offer
	CouponId       string    `dynamodbav:"couponId" json:"couponId"`
	UserId         string    `dynamodbav:"userId" json:"userId"`
	GeneratedAt    time.Time `dynamodbav:"generatedAt" json:"generatedAt"`
	ExpirationDate time.Time `dynamodbav:"validUntil" json:"validUntil"`
	Redeemed       bool      `dynamodbav:"redeemed" json:"redeemed"`
}

type CouponRange struct {
	Coupons []Coupon
	Next    *string
}

type OfferRange struct {
	Offers []GeneratedOffer
	Next   *string
}

func ValidatePassword(hashedPassword, plainTextPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword))
	return err == nil
}

func CreateTokenClient(c Client) string {
	now := time.Now()

	// valid for 6 hours
	validUntil := now.Add(time.Hour * 6).Unix()

	claims := jwt.MapClaims{
		"user":    c.Username,
		"email":   c.Email,
		"role":    "client",
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)

	// !STORE THE SECRET IN A MORE SECURE PLACE
	secret := []byte(SECRET)

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return ""
	}

	return tokenString

}

func CreateTokenEmployee(e Employee) string {
	now := time.Now()

	// valid for 6 hours
	validUntil := now.Add(time.Hour * 6).Unix()

	claims := jwt.MapClaims{
		"user":    e.Username,
		"email":   e.Email,
		"role":    "employee",
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)

	// !STORE THE SECRET IN A MORE SECURE PLACE
	secret := []byte(SECRET)

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return ""
	}

	return tokenString

}

func (u *User) HashPassword() error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)

	return nil
}
