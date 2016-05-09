package handler

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent_auth/command"
	"github.com/bborbe/bot_agent_auth/matcher"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type CreateApplication func(applicationName string) (*api.ApplicationPassword, error)

type handler struct {
	command           command.Command
	authToken         string
	createApplication CreateApplication
}

func New(prefix string, authToken string, createApplication CreateApplication) *handler {
	h := new(handler)
	h.command = command.New(prefix, "application", "create", "[NAME]")
	h.authToken = authToken
	h.createApplication = createApplication
	return h
}

func (h *handler) Match(request *message.Request) bool {
	return h.command.MatchRequest(request) && matcher.MatchRequestAuthToken(h.authToken, request)
}

func (h *handler) Help(request *message.Request) []string {
	if matcher.MatchRequestAuthToken(h.authToken, request) {
		return []string{h.command.Help()}
	}
	return []string{}
}

func (h *handler) HandleMessage(request *message.Request) ([]*message.Response, error) {
	applicationName, err := h.command.Parameter(request, "[NAME]")
	if err != nil {
		return nil, err
	}
	applicationPassword, err := h.createApplication(applicationName)
	if err != nil {
		logger.Debugf("application creation failed => send failure message: %v", err)
		return response.CreateReponseMessage(fmt.Sprintf("create application %s failed", applicationName)), nil
	}
	logger.Debugf("application created => send success message")
	return response.CreateReponseMessage(fmt.Sprintf("application %s created with password %s", applicationName, *applicationPassword)), nil
}
