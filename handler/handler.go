package handler

import (
	"github.com/bborbe/bot_agent/message_handler"
	"github.com/bborbe/bot_agent/message"
)

type Handler interface {
	Match(request *message.Request) bool
	message_handler.MessageHandler
	Help() string
}
