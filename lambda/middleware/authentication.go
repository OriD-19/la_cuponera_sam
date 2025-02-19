package middleware

import (
	"OriD19/webdev2/handlers"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

// middleware for validating JWT tokens and checking user permissions

// basic authentication header
func ValidateJWTMiddleware(ctx context.Context, next func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		tokenString := extractTokenFromHeaders(request.Headers)

		if strings.TrimSpace(tokenString) == "" {
			return handlers.ErrResponse(401, "missing JWT token"), nil
		}

		claims, err := parseToken(tokenString)

		if err != nil {
			return handlers.ErrResponse(http.StatusUnauthorized, err.Error()), nil
		}

		expires := int64(claims["expires"].(float64))

		if time.Now().Unix() > expires {
			return handlers.ErrResponse(http.StatusUnauthorized, "JWT token expired"), nil
		}

		return next(ctx, request)
	}
}

// client authorization header
func ValidateClientJWTMiddleware(ctx context.Context, next func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(c context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		tokenString := extractTokenFromHeaders(request.Headers)

		if strings.TrimSpace(tokenString) == "" {
			return handlers.ErrResponse(401, "missing JWT token"), nil
		}

		claims, err := parseToken(tokenString)

		if err != nil {
			return handlers.ErrResponse(http.StatusUnauthorized, err.Error()), nil
		}

		expires := int64(claims["expires"].(float64))

		if time.Now().Unix() > expires {
			return handlers.ErrResponse(http.StatusUnauthorized, "JWT token expired"), nil
		}

		role := claims["role"].(string)

		if role != "client" {
			return handlers.ErrResponse(http.StatusUnauthorized, "client role required"), nil
		}

		return next(c, request)
	}
}

// employee authorization header
func ValidateEmployeeJWTMiddleware(next func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(c context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		tokenString := extractTokenFromHeaders(request.Headers)

		if strings.TrimSpace(tokenString) == "" {
			return handlers.ErrResponse(401, "missing JWT token"), nil
		}

		claims, err := parseToken(tokenString)

		if err != nil {
			return handlers.ErrResponse(http.StatusUnauthorized, err.Error()), nil
		}

		expires := int64(claims["expires"].(float64))

		if time.Now().Unix() > expires {
			return handlers.ErrResponse(http.StatusUnauthorized, "JWT token expired"), nil
		}

		role := claims["role"].(string)

		if role != "employee" {
			return handlers.ErrResponse(http.StatusUnauthorized, "employee role required"), nil
		}

		return next(c, request)
	}
}

func extractTokenFromHeaders(headers map[string]string) string {
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

func parseToken(tokenString string) (jwt.MapClaims, error) {
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
