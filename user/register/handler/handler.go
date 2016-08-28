package handler

import (
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/response"
	"github.com/golang/glog"
)

type Register func(authToken api.AuthToken, userName string) error

type handler struct {
	command  command.Command
	register Register
}

func New(prefix string, register Register) *handler {
	h := new(handler)
	h.command = command.New(prefix, "register", "[NAME]")
	h.register = register
	return h
}

func (h *handler) Match(request *api.Request) bool {
	return h.command.MatchRequest(request)
}

func (h *handler) Help(request *api.Request) []string {
	return []string{h.command.Help()}
}

func (h *handler) HandleMessage(request *api.Request) ([]*api.Response, error) {
	glog.V(2).Infof("handle message: %s", request.Message)
	userName, err := h.command.Parameter(request, "[NAME]")
	if err != nil {
		return nil, err
	}
	glog.V(2).Infof("register user %s", userName)
	if err := h.register(request.AuthToken, userName); err != nil {
		glog.V(2).Infof("register %s failed: %v", userName, err)
		return response.CreateReponseMessage("register failed"), nil
	}
	glog.V(2).Infof("user %s  registered => send success message", userName)
	return response.CreateReponseMessage("registration completed"), nil
}
