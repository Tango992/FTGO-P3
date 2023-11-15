package repository

import "ungraded-3/models"

type Message interface {
	Post(*models.Message) error
	GetById(string) error
	GetBySenderReceiver(string, string) error
}