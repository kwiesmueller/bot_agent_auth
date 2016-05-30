package handler

import (
	"github.com/bborbe/bot_agent/api"
)

type Handler interface {
	Match(request *api.Request) bool
	api.MessageHandler
	Help(request *api.Request) []string
}
