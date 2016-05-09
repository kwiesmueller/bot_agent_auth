package handler

import (
	"fmt"
	"github.com/bborbe/auth/api"
	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
	"github.com/bborbe/bot_agent_auth/command"
)

var logger = log.DefaultLogger

type Whoami func(authToken string) (*api.UserName, error)

type handler struct {
	command command.Command
	whoami  Whoami
}

func New(prefix string, whoami Whoami) *handler {
	h := new(handler)
	h.command = command.New(prefix, "whoami")
	h.whoami = whoami
	return h
}

func (h *handler) Match(request *message.Request) bool {
	return h.command.MatchRequest(request)
}

func (h *handler) Help(request *message.Request) []string {
	return []string{h.command.Help()}
}

func (h *handler) HandleMessage(request *message.Request) ([]*message.Response, error) {
	userName, err := h.whoami(request.AuthToken)
	var name string
	if err != nil {
		logger.Debugf("get whoami failed: %v", err)
		name = "-"
	} else {
		name = string(*userName)
	}
	logger.Debugf("application whoamid => send success message")
	return response.CreateReponseMessage(
		fmt.Sprintf("UserToken %s", request.AuthToken),
		fmt.Sprintf("UserName: %s", name),
	), nil
}
