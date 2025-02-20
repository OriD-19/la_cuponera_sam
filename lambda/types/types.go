package types

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var SECRET = os.Getenv("SECRET")

const DATE_YYYY_MM_DD = "2006-01-02"

type Entity struct {
	EntityType string `json:"-" dynamodbav:"entityType" validator:"oneof=client administrator enterprise employee offer"`
}

// ************************************************************
// USER ENTITIES
// ************************************************************

type User struct {
	Entity
	Email     string    `dynamodbav:"email" json:"email" validator:"required,email"` // email is the id
	Username  string    `dynamodbav:"id" json:"username" validator:"required,max=100"`
	Password  string    `dynamodbav:"password" json:"-"  validator:"required,min=8"`
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
	FirstName    string `dynamodbav:"firstName" json:"firstName"`
	LastName     string `dynamodbav:"lastName" json:"lastName"`
	PhoneNumber  string `dynamodbav:"phoneNumber" json:"phoneNumber"`
	DUI          string `dynamodbav:"dui" json:"dui"`
	EnterpriseId string `dynamodbav:"enterpriseCode" json:"enterpriseCode"`
}

// ************************************************************
// COUPON ENTITIES
// ************************************************************

// store a YYYY_MM_DD time format, and parse it with an UnmarshalJSON custom method
type RegisteredCoupon struct {
	CouponId string `dynamodbav:"couponId" json:"couponId"`
	Status   string `dynamodbav:"status" json:"status"`
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
	EnterpriseId string `dynamodbav:"enterpriseCode" json:"enterpriseCode" validate:"required"`
	Category     string `dynamodbav:"category" json:"category"`
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
	Coupons []Coupon `json:"coupons"`
	Next    *string  `json:"next"`
}

type OfferRange struct {
	Offers []GeneratedOffer `json:"offers"`
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
		"username": c.Username,
		"email":    c.Email,
		"role":     "client",
		"expires":  validUntil,
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
		"username": e.Username,
		"email":    e.Email,
		"role":     "employee",
		"expires":  validUntil,
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

func HashPassword(password string) (string, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	password = string(hashedPassword)

	return password, nil
}

func ExtractTokenFromHeaders(headers map[string]string) string {
	authHeader, ok := headers["Authorization"]

	if !ok {
		return ""
	}

	splitToken := strings.Split(authHeader, "Bearer ")

	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		secret := []byte(os.Getenv("SECRET"))
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid JWT token")
	}

	// type assertion
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("failed to parse JWT claims")
	}

	return claims, nil
}

func GetClientAuthFromHeader(headers map[string]string) (Client, error) {
	tokenString := ExtractTokenFromHeaders(headers)

	if tokenString == "" {
		return Client{}, nil
	}

	claims, err := ParseToken(tokenString)

	if err != nil {
		return Client{}, err
	}

	username := claims["username"].(string)

	return Client{
		User: User{
			Username: username,
		},
	}, nil
}
