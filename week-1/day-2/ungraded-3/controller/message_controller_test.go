package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"ungraded-3/helpers"
	"ungraded-3/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	e *echo.Echo
)

func TestMain(m *testing.M) {
	e = echo.New()
	e.Validator = &helpers.CustomValidator{NewValidator: validator.New()}

	m.Run()
}

func TestPost(t *testing.T) {
	requestBody := strings.NewReader(`{
		"sender": "daniel@mail.com",
		"receiver": "john@mail.com",
		"type": "text",
		"content": "Pinjam dulu seratus"
	}`)

	message := &models.Message{
		Sender:   "daniel@mail.com",
		Receiver: "john@mail.com",
		Type:     "text",
		Content:  "Pinjam dulu seratus",
	}

	req := httptest.NewRequest(http.MethodPost, "localhost:8080/pesan", requestBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)

	mockRepository := NewMockRepository()
	mockRepository.Mock.On("Post", message).Return(message)

	messageController := NewMessageController(&mockRepository)
	messageController.Post(c)

	response := w.Result()
	body, _ := io.ReadAll(response.Body)
	
	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEmpty(t, responseBody)
	assert.Equal(t, responseBody["message"].(string), "Message posted")
}
