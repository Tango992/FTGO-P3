package helpers

import (
	"preview-w2/models"
	"preview-w2/utils"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func SignNewJWT(c echo.Context, user models.User) *utils.ErrResponse {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(2 * time.Hour).Unix(),
		"id": user.Id.Hex(),
		"email": user.Email,
		"name": user.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return utils.ErrInternalServer.New(err.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.Value = tokenString
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Expires = time.Now().Add(2 * time.Hour)
	c.SetCookie(cookie)

	return nil
}