package repository

import (
	"errors"
	"preview-w1/dto"
	"preview-w1/entity"
	"preview-w1/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type DbHandler struct {
	*gorm.DB
}

func NewDbHandler(db *gorm.DB) DbHandler {
	return DbHandler{
		DB: db,
	}
}

func (db DbHandler) AddUser(data *entity.User) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func (db DbHandler) FindUser(loginData dto.Login) (entity.User, error) {
	var user entity.User

	res := db.Where("email = ?", loginData.Email).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return entity.User{}, echo.NewHTTPError(utils.ErrUnauthorized.Details("Invalid email / password"))
		}
		return entity.User{}, echo.NewHTTPError(utils.ErrInternalServer.Details(res.Error.Error()))
	}
	return user, nil
}

func (db DbHandler) AddLoan(data *entity.Loan) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func (db DbHandler) CheckLimit(userId uint) (float32, error) {
	var data float32
	
	if err := db.Model(&entity.Loan{}).Select("loan").Where("user_id = ?", userId).Take(&data).Error; err != nil {
		return 0, err
	}
	return data, nil
}