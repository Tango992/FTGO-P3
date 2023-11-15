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
	messageId := c.Param("id")

	message, err := mc.Repository.GetById(messageId)
	if err != nil {
		return c.JSON(err.Format())
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Get message by id",
		"data": message,
	})
}

func (mc MessageController) GetBySender(c echo.Context) error {
	senderEmail := c.Param("senderEmail")

	messages, err := mc.Repository.GetBySenderReceiver("sender", senderEmail)
	if err != nil {
		return c.JSON(err.Format())
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Get message by sender",
		"data": messages,
	})
}

func (mc MessageController) GetByReceiver(c echo.Context) error {
	receiverEmail := c.Param("receiverEmail")

	messages, err := mc.Repository.GetBySenderReceiver("receiver", receiverEmail)
	if err != nil {
		return c.JSON(err.Format())
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Get message by receiver",
		"data": messages,
	})
}