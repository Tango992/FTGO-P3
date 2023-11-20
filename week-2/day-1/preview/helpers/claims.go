package helpers

import (
	"preview-w2/dto"
	"preview-w2/utils"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GetClaims(c echo.Context) (dto.Claims, error) {
	claimsTmp := c.Get("user")
	if claimsTmp == nil {
		return dto.Claims{}, echo.NewHTTPError(utils.ErrUnauthorized.EchoFormatDetails("Failed to fetch user claims from JWT"))
	}
	
	claims := claimsTmp.(jwt.MapClaims)
	return dto.Claims{
		ID:       claims["id"].(string),
		Email:    claims["email"].(string),
		Name: claims["name"].(string),
	}, nil
}
