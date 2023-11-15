package controller

import (
	"net/http"
	"ungraded-3/dto"
	"ungraded-3/models"
	"ungraded-3/repository"
	"ungraded-3/utils"

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
	var postMessageTmp dto.PostMessage
	if err := c.Bind(&postMessageTmp); err != nil {
		return c.JSON(utils.ErrBadRequest.Details(err.Error()))
	}

	if err := c.Validate(&postMessageTmp); err != nil {
		return c.JSON(utils.ErrBadRequest.Details(err.Error()))
	}
	
	postMessage := models.Message{
		Sender: postMessageTmp.Sender,
		Receiver: postMessageTmp.Receiver,
		Type: postMessageTmp.Type,
		Content: postMessageTmp.Content,
	}
	
	if err := mc.Repository.Post(&postMessage); err != nil {
		return c.JSON(err.Format())
	}
	
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Message posted",
		"data": postMessage,
	})
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