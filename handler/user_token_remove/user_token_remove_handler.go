package user_token_remove

import (
	auth_model "github.com/bborbe/auth/model"

	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/golang/glog"
)

type removeTokenToUser func(newToken auth_model.AuthToken, username auth_model.UserName) error

type handler struct {
	command           command.Command
	removeTokenToUser removeTokenToUser
}

func New(
	prefix model.Prefix,
	removeTokenToUser removeTokenToUser,
) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "remove", "token", "[TOKEN]", "from", "user", "[USERNAME]")
	h.removeTokenToUser = removeTokenToUser
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
	username, err := h.command.Parameter(request, "[USERNAME]")
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("remove token %v to user %v", token, username)
	if err := h.removeTokenToUser(auth_model.AuthToken(token), auth_model.UserName(username)); err != nil {
		glog.V(4).Infof("remove token %v to user %v failed: %v", token, username, err)
		return response.CreateReponseMessage("remove token failed"), nil
	}
	glog.V(2).Infof("token %v removeed to user %v", token, username)
	return response.CreateReponseMessage("token removed to user"), nil
}
