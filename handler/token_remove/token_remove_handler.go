package handler

import (
	auth_model "github.com/bborbe/auth/model"

	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/golang/glog"
)

type remove func(authToken auth_model.AuthToken, token auth_model.AuthToken) error

type handler struct {
	command command.Command
	remove  remove
}

func New(prefix model.Prefix, remove remove) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "token", "remove", "[TOKEN]")
	h.remove = remove
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
	token, err := h.command.Parameter(request, "[TOKEN]")
	if err != nil {
		return nil, err
	}
	glog.V(2).Infof("remove token %s", token)
	if err := h.remove(request.AuthToken, auth_model.AuthToken(token)); err != nil {
		glog.V(2).Infof("remove token %s failed: %v", token, err)
		return response.CreateReponseMessage("remove token failed"), nil
	}
	glog.V(2).Infof("token %s removed => send success message", token)
	return response.CreateReponseMessage("token removed"), nil
}
