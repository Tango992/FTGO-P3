package models

type User struct {
	Id         uint    `json:"id" bson:"_id,omitempty" validate:"required"`
	Name       string  `json:"name" validate:"required"`
	Address    string  `json:"address" validate:"required"`
	Email      string  `json:"email" validate:"required,email"`
	Salary     float32 `json:"salary" validate:"required"`
	Department string  `json:"department" validate:"required"`
}
