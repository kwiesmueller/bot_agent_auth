package user_create

import (
	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/bborbe/http/header"
	"github.com/golang/glog"
)

type create func(userName auth_model.UserName, authToken auth_model.AuthToken) error

type handler struct {
	command command.Command
	create  create
}

func New(
	prefix model.Prefix,
	create create,
) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "user", "create", "[USERNAME]", "[PASSWORD]")
	h.create = create
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
	userName, err := h.command.Parameter(request, "[USERNAME]")
	if err != nil {
		return nil, err
	}
	password, err := h.command.Parameter(request, "[PASSWORD]")
	if err != nil {
		return nil, err
	}
	authToken := auth_model.AuthToken(header.CreateAuthorizationToken(userName, password))
	glog.V(2).Infof("create user %s", userName)
	if err := h.create(auth_model.UserName(userName), authToken); err != nil {
		glog.V(2).Infof("create %s failed: %v", userName, err)
		return response.CreateReponseMessage("create failed"), nil
	}
	glog.V(2).Infof("user %s  created => send success message", userName)
	return response.CreateReponseMessage("create user completed"), nil
}
