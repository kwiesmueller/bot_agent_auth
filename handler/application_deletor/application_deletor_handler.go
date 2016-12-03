package application_deletor

import (
	"fmt"
	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/golang/glog"
)

type deleteApplication func(applicationName auth_model.ApplicationName) error

type handler struct {
	command           command.Command
	deleteApplication deleteApplication
}

func New(
	prefix model.Prefix,
	deleteApplication deleteApplication,
) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "application", "delete", "[NAME]")
	h.deleteApplication = deleteApplication
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
	glog.V(2).Infof("delete applcation %s", applicationName)
	if err := h.deleteApplication(auth_model.ApplicationName(applicationName)); err != nil {
		return response.CreateReponseMessage(fmt.Sprintf("delete application %s failed", applicationName)), nil
	}
	return response.CreateReponseMessage(fmt.Sprintf("application %s deleted", applicationName)), nil
}
