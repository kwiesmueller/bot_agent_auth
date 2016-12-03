package user_whoami

import (
	"fmt"

	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/golang/glog"
)

type whoami func(authToken auth_model.AuthToken) (*auth_model.UserName, error)

type handler struct {
	command command.Command
	whoami  whoami
}

func New(prefix model.Prefix, whoami whoami) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "whoami")
	h.whoami = whoami
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
	userName, err := h.whoami(request.AuthToken)
	var name string
	if err != nil {
		glog.V(2).Infof("get whoami failed: %v", err)
		name = "-"
	} else {
		name = string(*userName)
	}
	glog.V(2).Infof("application whoamid => send success message")
	return response.CreateReponseMessage(
		fmt.Sprintf("UserToken %s", request.AuthToken),
		fmt.Sprintf("UserName: %s", name),
	), nil
}
