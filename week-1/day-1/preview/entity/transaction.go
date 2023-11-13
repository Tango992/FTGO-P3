package entity

import "time"

type Transaction struct {
	ID          uint    `gorm:"primaryKey"`
	UserID      uint    `gorm:"not null"`
	Transaction string  `gorm:"not null"`
	Amount      float32 `gorm:"not null;check:(amount >= 0)"`
	CreatedAt   time.Time
}
