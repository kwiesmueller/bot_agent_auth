package handler

import (
	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
	"github.com/bborbe/bot_agent_auth/command"
)

var logger = log.DefaultLogger

type Register func(authToken string, userName string) error

type handler struct {
	command  command.Command
	register Register
}

func New(prefix string, register Register) *handler {
	h := new(handler)
	h.command = command.New(prefix, "register", "[NAME]")
	h.register = register
	return h
}

func (h *handler) Match(request *message.Request) bool {
	return h.command.MatchRequest(request)
}

func (h *handler) Help(request *message.Request) []string {
	return []string{h.command.Help()}
}

func (h *handler) HandleMessage(request *message.Request) ([]*message.Response, error) {
	logger.Debugf("handle message: %s", request.Message)
	userName, err := h.command.Parameter(request, "[NAME]")
	if err != nil {
		return nil, err
	}
	logger.Debugf("register user %s", userName)
	if err := h.register(request.AuthToken, userName); err != nil {
		logger.Debugf("register %s failed: %v", userName, err)
		return response.CreateReponseMessage("register failed"), nil
	}
	logger.Debugf("user %s  registered => send success message", userName)
	return response.CreateReponseMessage("registration completed"), nil
}
