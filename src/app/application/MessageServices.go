package application

import (
	"awesomeProject1/src/app/domain"
	"awesomeProject1/src/app/infrastructure/Repositories"
	"github.com/gorilla/websocket"
)

type ViewMessage struct {
	Message domain.Message `json:"message"`
	User    domain.User    `json:"user"`
}

var messageRepo = Repositories.NewMessageRepo()

type (
	MessageService struct {
	}
)

func NewMessageService() MessageService {
	return MessageService{}
}

func (r MessageService) CreateMessage(text string, conn *websocket.Conn, userId string) bool {
	message, _ := messageRepo.CreateMessage(text, userId)

	var viewMessage ViewMessage

	viewMessage.Message = message
	user, _ := Repositories.NewUserRepo().UserByID(message.UserId)
	viewMessage.User = user

	if conn != nil {
		conn.WriteJSON(viewMessage)
	}

	return true
}

func (r MessageService) GetLastMessages() []ViewMessage {
	var viewMessages []ViewMessage

	messages := messageRepo.GetFiveLastMessages()
	for _, s := range messages {
		var v ViewMessage
		v.Message = s
		user, _ := Repositories.NewUserRepo().UserByID(s.UserId)
		v.User = user
		viewMessages = append(viewMessages, v)
	}

	return viewMessages
}
