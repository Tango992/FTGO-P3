package repository

import (
	"ungraded-3/models"
	"ungraded-3/utils"
)

type Message interface {
	Post(*models.Message) *utils.ErrResponse
	GetById(string) *utils.ErrResponse
	GetBySenderReceiver(string, string) *utils.ErrResponse
}