package handler

import (
	"fmt"

	auth_model "github.com/bborbe/auth/model"

	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/matcher"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/golang/glog"
)

type DeleteApplication func(applicationName string) error

type handler struct {
	command           command.Command
	authToken         auth_model.AuthToken
	deleteApplication DeleteApplication
}

func New(prefix model.Prefix, authToken auth_model.AuthToken, deleteApplication DeleteApplication) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "application", "delete", "[NAME]")
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
	glog.V(2).Infof("delete applcation %s", applicationName)
	if err := h.deleteApplication(applicationName); err != nil {
		return response.CreateReponseMessage(fmt.Sprintf("delete application %s failed", applicationName)), nil
	}
	return response.CreateReponseMessage(fmt.Sprintf("application %s deleted", applicationName)), nil
}
