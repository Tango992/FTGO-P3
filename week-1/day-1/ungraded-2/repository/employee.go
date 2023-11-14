package repository

import (
	"ungraded-2/dto"
	"ungraded-2/models"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Employee interface {
	Create(*models.User) error
	Find(*options.FindOptions) ([]models.User, error)
	FindById(string) (models.User, error)
	Update(dto.UpdateEmployee, string) error
	DeleteById(string) error
}