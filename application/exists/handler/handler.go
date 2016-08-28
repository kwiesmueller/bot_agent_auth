package handler

import (
	"fmt"

	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/matcher"
	"github.com/bborbe/bot_agent/response"
	"github.com/golang/glog"
)

type ExistsApplication func(applicationName string) (bool, error)

type handler struct {
	command           command.Command
	authToken         api.AuthToken
	existsApplication ExistsApplication
}

func New(prefix string, authToken api.AuthToken, existsApplication ExistsApplication) *handler {
	h := new(handler)
	h.command = command.New(prefix, "application", "exists", "[NAME]")
	h.authToken = authToken
	h.existsApplication = existsApplication
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
	exists, err := h.existsApplication(applicationName)
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
