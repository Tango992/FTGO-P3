package entity

import "time"

type Loan struct {
	UserID    uint    `gorm:"primaryKey"`
	Salary    float32 `gorm:"not null"`
	Loan      float32 `gorm:"not null;check:(loan >= 0)"`
	CreatedAt time.Time
}
