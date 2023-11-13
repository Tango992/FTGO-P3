package dto

type Loan struct {
	Salary float32 `json:"salary" validate:"required"`
	Loan   float32 `json:"loan" validate:"required"`
}
