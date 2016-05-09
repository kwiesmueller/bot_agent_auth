package handler

import (
	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent_auth/command"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type Register func(authToken string, userName string) error

type handler struct {
	command  command.Command
	register Register
}

func New(prefix string, register Register) *handler {
	h := new(handler)
	h.command = command.New(prefix, "token", "add", "[NAME]")
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
	token, err := h.command.Parameter(request, "[NAME]")
	if err != nil {
		return nil, err
	}
	logger.Debugf("add token %s", token)
	if err := h.register(request.AuthToken, token); err != nil {
		logger.Debugf("add token %s failed: %v", token, err)
		return response.CreateReponseMessage("add token failed"), nil
	}
	logger.Debugf("token %s added => send success message", token)
	return response.CreateReponseMessage("token added"), nil
}
