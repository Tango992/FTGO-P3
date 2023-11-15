package repository

import (
	"ungraded-3/models"
	"ungraded-3/utils"
)

type Message interface {
	Post(*models.Message) *utils.ErrResponse
	GetById(string) (models.Message, *utils.ErrResponse)
	GetBySenderReceiver(string, string) ([]models.Message, *utils.ErrResponse)
}