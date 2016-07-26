package handler

import (
	"fmt"

	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type Whoami func(authToken string) (*auth_model.UserName, error)

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

func (h *handler) Match(request *api.Request) bool {
	return h.command.MatchRequest(request)
}

func (h *handler) Help(request *api.Request) []string {
	return []string{h.command.Help()}
}

func (h *handler) HandleMessage(request *api.Request) ([]*api.Response, error) {
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
