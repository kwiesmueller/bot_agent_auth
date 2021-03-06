package user_register

import (
	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/golang/glog"
)

type register func(userName auth_model.UserName, authToken auth_model.AuthToken) error

type handler struct {
	command  command.Command
	register register
}

func New(prefix model.Prefix, register register) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "register", "[NAME]")
	h.register = register
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
	glog.V(2).Infof("handle message: %s", request.Message)
	userName, err := h.command.Parameter(request, "[NAME]")
	if err != nil {
		return nil, err
	}
	glog.V(2).Infof("register user %s", userName)
	if err := h.register(auth_model.UserName(userName), request.AuthToken); err != nil {
		glog.V(2).Infof("register %s failed: %v", userName, err)
		return response.CreateReponseMessage("register failed"), nil
	}
	glog.V(2).Infof("user %s  registered => send success message", userName)
	return response.CreateReponseMessage("registration completed"), nil
}
