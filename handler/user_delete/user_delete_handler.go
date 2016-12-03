package user_delete

import (
	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/golang/glog"
)

type deleteUser func(username auth_model.UserName) error

type handler struct {
	command command.Command
	delete  deleteUser
}

func New(
	prefix model.Prefix,
	delete deleteUser,
) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "user", "delete", "[USERNAME]")
	h.delete = delete
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
	glog.V(2).Infof("delete user")
	username, err := h.command.Parameter(request, "[USERNAME]")
	if err != nil {
		return nil, err
	}
	glog.V(2).Infof("delete user %s", username)
	if err := h.delete(auth_model.UserName(username)); err != nil {
		glog.V(2).Infof("delete user %s failed: %v", username, err)
		return response.CreateReponseMessage("delete failed"), nil
	}
	glog.V(2).Infof("delete user %s successful", username)
	return response.CreateReponseMessage("delete user completed"), nil
}
