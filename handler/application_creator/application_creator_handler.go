package application_creator

import (
	"fmt"
	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/golang/glog"
)

type createApplication func(applicationName auth_model.ApplicationName) (*auth_model.ApplicationPassword, error)

type handler struct {
	command           command.Command
	createApplication createApplication
}

func New(
	prefix model.Prefix,
	createApplication createApplication,
) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "application", "create", "[NAME]")
	h.createApplication = createApplication
	return h
}

func (h *handler) Allowed(request *api.Request) bool {
	return true
}

func (h *handler) Match(request *api.Request) bool {
	return h.command.MatchRequest(request)
}

func (h *handler) Help(request *api.Request) []string {
	return []string{h.command.Help()}
}

func (h *handler) HandleMessage(request *api.Request) ([]*api.Response, error) {
	applicationName, err := h.command.Parameter(request, "[NAME]")
	if err != nil {
		return nil, err
	}
	applicationPassword, err := h.createApplication(auth_model.ApplicationName(applicationName))
	if err != nil {
		glog.V(2).Infof("application creation failed => send failure message: %v", err)
		return response.CreateReponseMessage(fmt.Sprintf("create application %s failed", applicationName)), nil
	}
	glog.V(3).Infof("application created => send success message")
	return response.CreateReponseMessage(fmt.Sprintf("application %s created with password %s", applicationName, *applicationPassword)), nil
}
