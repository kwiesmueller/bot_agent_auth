package handler

import (
	"fmt"

	auth_api "github.com/bborbe/auth/model"
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/matcher"
	"github.com/bborbe/bot_agent/response"
	"github.com/golang/glog"
)

type CreateApplication func(applicationName string) (*auth_api.ApplicationPassword, error)

type handler struct {
	command           command.Command
	authToken         api.AuthToken
	createApplication CreateApplication
}

func New(prefix string, authToken api.AuthToken, createApplication CreateApplication) *handler {
	h := new(handler)
	h.command = command.New(prefix, "application", "create", "[NAME]")
	h.authToken = authToken
	h.createApplication = createApplication
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
	applicationPassword, err := h.createApplication(applicationName)
	if err != nil {
		glog.V(2).Infof("application creation failed => send failure message: %v", err)
		return response.CreateReponseMessage(fmt.Sprintf("create application %s failed", applicationName)), nil
	}
	glog.V(2).Infof("application created => send success message")
	return response.CreateReponseMessage(fmt.Sprintf("application %s created with password %s", applicationName, *applicationPassword)), nil
}
