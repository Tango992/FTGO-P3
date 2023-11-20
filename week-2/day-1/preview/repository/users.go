package repository

import (
	"preview-w2/dto"
	"preview-w2/models"
	"preview-w2/utils"
)

type Users interface {
	Register(*models.User) *utils.ErrResponse
	FindUser(dto.LoginUser) (models.User, *utils.ErrResponse)
}