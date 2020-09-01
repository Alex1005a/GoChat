package application

import (

	"awesomeProject1/src/app/infrastructure/Repositories"
	"github.com/gorilla/websocket"
)

var messageRepo = Repositories.NewMessageRepo()

type (
	MessageService struct {
	}
)

func NewMessageService() MessageService {
	return MessageService{}
}

func (r MessageService) CreateMessage(text string, conn *websocket.Conn) bool {
	message, _ := messageRepo.CreateMessage(text)
	conn.WriteJSON(message)

	return true
}