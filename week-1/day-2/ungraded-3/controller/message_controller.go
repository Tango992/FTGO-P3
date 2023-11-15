package controller

import (
	"ungraded-3/repository"

	"github.com/labstack/echo/v4"
)

type MessageController struct {
	Repository repository.Message
}

func NewMessageController(r repository.Message) MessageController {
	return MessageController{
		Repository: r,
	}
}

func (mc MessageController) Post(c echo.Context) error {
	return nil
}

func (mc MessageController) GetById(c echo.Context) error {
	return nil
}

func (mc MessageController) GetBySender(c echo.Context) error {
	return nil
}

func (mc MessageController) GetByReceiver(c echo.Context) error {
	return nil
}