package entity

type User struct {
	ID           uint   `gorm:"primaryKey"`
	FullName     string `gorm:"not null"`
	Email        string `gorm:"not null"`
	Password     string `gorm:"not null"`
	Birth        string `gorm:"not null"`
	Transactions []Transaction
	Loan         Loan
}
