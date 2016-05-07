package handler

import (
	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent/message_handler"
)

type Handler interface {
	Match(request *message.Request) bool
	message_handler.MessageHandler
	Help() string
}
