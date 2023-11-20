package dto

type Payment struct {
	Amount  float32 `json:"amount" validate:"required"`
	Channel string  `json:"channel" validate:"required,oneof=akulaku shopeepaylater"`
}
