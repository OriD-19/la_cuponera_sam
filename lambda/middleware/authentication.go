package middleware

import (
	"OriD19/webdev2/handlers"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

// middleware for validating JWT tokens and checking user permissions

// basic authentication header
func ValidateJWTMiddleware(next func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		tokenString := extractTokenFromHeaders(request.Headers)

		if strings.TrimSpace(tokenString) == "" {
			return handlers.Response(401, "missing JWT token"), nil
		}

		claims, err := parseToken(tokenString)

		if err != nil {
			return handlers.Response(http.StatusUnauthorized, err.Error()), nil
		}

		expires := int64(claims["expires"].(float64))

		if time.Now().Unix() > expires {
			return handlers.Response(http.StatusUnauthorized, "JWT token expired"), nil
		}

		return next(request)
	}
}

// client authorization header
func ValidateClientJWTMiddleware(next func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		tokenString := extractTokenFromHeaders(request.Headers)

		if strings.TrimSpace(tokenString) == "" {
			return handlers.Response(401, "missing JWT token"), nil
		}

		claims, err := parseToken(tokenString)

		if err != nil {
			return handlers.Response(http.StatusUnauthorized, err.Error()), nil
		}

		expires := int64(claims["expires"].(float64))

		if time.Now().Unix() > expires {
			return handlers.Response(http.StatusUnauthorized, "JWT token expired"), nil
		}

		role := claims["role"].(string)

		if role != "client" {
			return handlers.Response(http.StatusUnauthorized, "client role required"), nil
		}

		return next(request)
	}
}

// employee authorization header
func ValidateEmployeeJWTMiddleware(next func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		tokenString := extractTokenFromHeaders(request.Headers)

		if strings.TrimSpace(tokenString) == "" {
			return handlers.Response(401, "missing JWT token"), nil
		}

		claims, err := parseToken(tokenString)

		if err != nil {
			return handlers.Response(http.StatusUnauthorized, err.Error()), nil
		}

		expires := int64(claims["expires"].(float64))

		if time.Now().Unix() > expires {
			return handlers.Response(http.StatusUnauthorized, "JWT token expired"), nil
		}

		role := claims["role"].(string)

		if role != "employee" {
			return handlers.Response(http.StatusUnauthorized, "employee role required"), nil
		}

		return next(request)
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
		secret := []byte("secret")
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
