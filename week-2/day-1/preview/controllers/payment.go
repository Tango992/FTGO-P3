package controllers

import (
	"net/http"
	"preview-w2/dto"
	"preview-w2/helpers"
	"preview-w2/models"
	"preview-w2/repository"
	"preview-w2/utils"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentController struct {
	Repository repository.Payments
}

func NewPaymentController(repository repository.Payments) PaymentController {
	return PaymentController{
		Repository: repository,
	}
}

func (p PaymentController) CreatePayment(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}
	userId, _ := primitive.ObjectIDFromHex(user.ID)
	
	var paymentDataTmp dto.Payment
	if err := c.Bind(&paymentDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}
	
	if err := c.Validate(&paymentDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	paymentData := models.Payment{
		UserId: userId,
		Amount: paymentDataTmp.Amount,
		Channel: paymentDataTmp.Channel,
		Status: "payment",
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	if err := p.Repository.Create(&paymentData); err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}
	
	return c.JSON(http.StatusCreated, dto.Response{
		Message: "Payment created",
		Data: paymentData,
	})
}