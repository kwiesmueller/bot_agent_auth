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

type ListUsers func() ([]auth_model.UserName, error)

type handler struct {
	command   command.Command
	authToken auth_model.AuthToken
	listUsers ListUsers
}

func New(prefix model.Prefix, authToken auth_model.AuthToken, list ListUsers) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "user", "list")
	h.authToken = authToken
	h.listUsers = list
	return h
}

func (h *handler) Match(request *api.Request) bool {
	return h.command.MatchRequest(request) && matcher.MatchRequestAuthToken(h.authToken, request)
}

func (h *handler) Help(request *api.Request) []string {
	return []string{h.command.Help()}
}

func (h *handler) HandleMessage(request *api.Request) ([]*api.Response, error) {
	userNames, err := h.listUsers()
	glog.V(2).Infof("user list => send success message")
	if err != nil {
		glog.V(2).Infof("list user failed: %v", err)
		return response.CreateReponseMessage("list user failed"), nil
	}
	var results []string
	for _, userName := range userNames {
		results = append(results, string(userName))
	}
	return response.CreateReponseMessage(results...), nil
}
