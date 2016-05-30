package handler

import (
	"fmt"

	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent_auth/command"
	"github.com/bborbe/bot_agent_auth/matcher"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type DeleteApplication func(applicationName string) error

type handler struct {
	command           command.Command
	authToken         string
	deleteApplication DeleteApplication
}

func New(prefix string, authToken string, deleteApplication DeleteApplication) *handler {
	h := new(handler)
	h.command = command.New(prefix, "application", "delete", "[NAME]")
	h.authToken = authToken
	h.deleteApplication = deleteApplication
	return h
}

func (h *handler) Match(request *api.Request) bool {
	return h.command.MatchRequest(request) && matcher.MatchRequestAuthToken(h.authToken, request)
}

func (h *handler) Help(request *api.Request) []string {
	if matcher.MatchRequestAuthToken(h.authToken, request) {
		return []string{h.command.Help()}
	}
	return []string{}
}

func (h *handler) HandleMessage(request *api.Request) ([]*api.Response, error) {
	applicationName, err := h.command.Parameter(request, "[NAME]")
	if err != nil {
		return nil, err
	}
	logger.Debugf("delete applcation %s", applicationName)
	if err := h.deleteApplication(applicationName); err != nil {
		return response.CreateReponseMessage(fmt.Sprintf("delete application %s failed", applicationName)), nil
	}
	return response.CreateReponseMessage(fmt.Sprintf("application %s deleted", applicationName)), nil
}
