package repository

import (
	"preview-w2/models"
	"preview-w2/utils"
)

type Payments interface {
	Create(*models.Payment) *utils.ErrResponse
}