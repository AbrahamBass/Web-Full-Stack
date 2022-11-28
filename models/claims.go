package models

import "github.com/golang-jwt/jwt/v4"

type AppClient struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}
