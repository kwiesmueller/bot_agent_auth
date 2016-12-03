package user_token_add

import (
	auth_model "github.com/bborbe/auth/model"

	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/golang/glog"
)

type addTokenToUser func(newToken auth_model.AuthToken, username auth_model.UserName) error

type handler struct {
	command        command.Command
	addTokenToUser addTokenToUser
}

func New(
	prefix model.Prefix,
	addTokenToUser addTokenToUser,
) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "user", "[USERNAME]", "add", "token", "[TOKEN]")
	h.addTokenToUser = addTokenToUser
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
	token, err := h.command.Parameter(request, "[TOKEN]")
	if err != nil {
		return nil, err
	}
	username, err := h.command.Parameter(request, "[USERNAME]")
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("add token %v to user %v", token, username)
	if err := h.addTokenToUser(auth_model.AuthToken(token), auth_model.UserName(username)); err != nil {
		glog.V(4).Infof("add token %v to user %v failed: %v", token, username, err)
		return response.CreateReponseMessage("add token failed"), nil
	}
	glog.V(2).Infof("token %v added to user %v", token, username)
	return response.CreateReponseMessage("token added to user"), nil
}
