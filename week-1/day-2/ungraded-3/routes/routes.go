package routes

import (
	"ungraded-3/controller"

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo, messageController controller.MessageController) {
	e.POST("/pesan", messageController.Post)
	e.GET("/pesan/:id", messageController.GetById)
	e.GET("/pesan/pengguna/:senderEmail", messageController.GetBySender)
	e.GET("/pesan/penerima/:receiverEmail", messageController.GetByReceiver)
}
