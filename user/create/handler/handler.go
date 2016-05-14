package handler

import (
	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent_auth/command"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
	"github.com/bborbe/http/header"
	"github.com/bborbe/auth/api"
)

var logger = log.DefaultLogger

type Create func(userName string, authToken string) error

type handler struct {
	command command.Command
	create  Create
}

func New(prefix string, create Create) *handler {
	h := new(handler)
	h.command = command.New(prefix, "create", "[USERNAME]", "[PASSWORD]")
	h.create = create
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
	userName, err := h.command.Parameter(request, "[USERNAME]")
	if err != nil {
		return nil, err
	}
	password, err := h.command.Parameter(request, "[PASSWORD]")
	if err != nil {
		return nil, err
	}
	authToken := header.CreateAuthorizationToken(userName, password)
	logger.Debugf("create user %s", userName)
	if err := h.create(userName, authToken); err != nil {
		logger.Debugf("create %s failed: %v", userName, err)
		return response.CreateReponseMessage("create failed"), nil
	}
	logger.Debugf("user %s  created => send success message", userName)
	return response.CreateReponseMessage("create user completed"), nil
}
