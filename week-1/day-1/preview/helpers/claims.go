package helpers

import (
	"preview-w1/dto"
	"preview-w1/utils"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GetClaims(c echo.Context) (dto.Claims, error) {
	claimsTmp := c.Get("user")
	if claimsTmp == nil {
		return dto.Claims{}, echo.NewHTTPError(utils.ErrUnauthorized.Details("Failed to fetch user claims from JWT"))
	}
	
	claims := claimsTmp.(jwt.MapClaims)
	return dto.Claims{
		ID:       uint(claims["id"].(float64)),
		Email:    claims["email"].(string),
		FullName: claims["full_name"].(string),
	}, nil
}
