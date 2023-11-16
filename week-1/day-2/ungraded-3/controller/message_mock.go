package controller

import (
	"ungraded-3/models"
	"ungraded-3/utils"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockMessageRepository struct {
	Mock mock.Mock
}

func NewMockRepository() MockMessageRepository {
	return MockMessageRepository{}
}

func (m *MockMessageRepository) Post(data *models.Message) *utils.ErrResponse {
	args := m.Mock.Called(data)
	args[0].(*models.Message).ID = primitive.NilObjectID
	return nil
}

func (m *MockMessageRepository) GetById(msgId string) (models.Message, *utils.ErrResponse) {
	args := m.Mock.Called(msgId)
	return args.Get(0).(models.Message), args.Get(1).(*utils.ErrResponse)
}

func (m *MockMessageRepository) GetBySenderReceiver(subject, email string) ([]models.Message, *utils.ErrResponse) {
	args := m.Mock.Called(subject, email)
	return args.Get(0).([]models.Message), args.Get(1).(*utils.ErrResponse)
}
