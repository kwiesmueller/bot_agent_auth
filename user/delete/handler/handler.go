package handler

import (
	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/matcher"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/golang/glog"
)

type DeleteUser func(username string) error

type handler struct {
	command   command.Command
	authToken auth_model.AuthToken
	delete    DeleteUser
}

func New(prefix model.Prefix, authToken auth_model.AuthToken, delete DeleteUser) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "user", "delete", "[USERNAME]")
	h.authToken = authToken
	h.delete = delete
	return h
}

func (h *handler) Match(request *api.Request) bool {
	return h.command.MatchRequest(request) && matcher.MatchRequestAuthToken(h.authToken, request)
}

func (h *handler) Help(request *api.Request) []string {
	return []string{h.command.Help()}
}

func (h *handler) HandleMessage(request *api.Request) ([]*api.Response, error) {
	glog.V(2).Infof("delete user")
	username, err := h.command.Parameter(request, "[USERNAME]")
	if err != nil {
		return nil, err
	}
	glog.V(2).Infof("delete user %s", username)
	if err := h.delete(username); err != nil {
		glog.V(2).Infof("delete user %s failed: %v", username, err)
		return response.CreateReponseMessage("delete failed"), nil
	}
	glog.V(2).Infof("delete user %s successful", username)
	return response.CreateReponseMessage("delete user completed"), nil
}
