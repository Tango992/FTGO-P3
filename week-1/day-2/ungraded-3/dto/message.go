package dto

type PostMessage struct {
	Sender   string             `json:"sender" bson:"sender" validate:"required,email"`
	Receiver string             `json:"receiver" bson:"receiver" validate:"required,email"`
	Type     string             `json:"type" bson:"type" validate:"required,oneof=text image file"`
	Content  string             `json:"content" bson:"content" validate:"required"`
}
