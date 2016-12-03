package application_exists

import (
	"fmt"

	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/golang/glog"
)

type existsApplication func(applicationName auth_model.ApplicationName) (bool, error)

type handler struct {
	command           command.Command
	existsApplication existsApplication
}

func New(
	prefix model.Prefix,
	existsApplication existsApplication,
) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "application", "exists", "[NAME]")
	h.existsApplication = existsApplication
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
	exists, err := h.existsApplication(auth_model.ApplicationName(applicationName))
	if err != nil {
		glog.V(2).Infof("application exists failed => send failure message: %v", err)
		return response.CreateReponseMessage(fmt.Sprintf("exists application %s failed", applicationName)), nil
	}
	glog.V(2).Infof("application exists => send success message")
	if exists {
		return response.CreateReponseMessage(fmt.Sprintf("application %s exists", applicationName)), nil
	} else {
		return response.CreateReponseMessage(fmt.Sprintf("application %s not exists", applicationName)), nil
	}
}
