package handler

import (
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent_auth/command"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type Unregister func(authToken string) error

type handler struct {
	command    command.Command
	unregister Unregister
}

func New(prefix string, unregister Unregister) *handler {
	h := new(handler)
	h.command = command.New(prefix, "unregister")
	h.unregister = unregister
	return h
}

func (h *handler) Match(request *api.Request) bool {
	return h.command.MatchRequest(request)
}

func (h *handler) Help(request *api.Request) []string {
	return []string{h.command.Help()}
}

func (h *handler) HandleMessage(request *api.Request) ([]*api.Response, error) {
	logger.Debugf("unregister user with token %s", request.AuthToken)
	if err := h.unregister(request.AuthToken); err != nil {
		logger.Debugf("unregister user with token %s failed: %v", request.AuthToken, err)
		return response.CreateReponseMessage("unregister failed"), nil
	}
	logger.Debugf("unregister user with token %s successful", request.AuthToken)
	return response.CreateReponseMessage("unregistration completed"), nil
}
