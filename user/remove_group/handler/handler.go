package handler

import (
	"fmt"
	auth_model "github.com/bborbe/auth/model"

	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/matcher"
	"github.com/bborbe/bot_agent/response"
	"github.com/golang/glog"
)

type RemoveGroupToUser func(groupName string, userName string) error

type handler struct {
	command           command.Command
	authToken         auth_model.AuthToken
	removeGroupToUser RemoveGroupToUser
}

func New(prefix string, authToken auth_model.AuthToken, removeGroupToUser RemoveGroupToUser) *handler {
	h := new(handler)
	h.command = command.New(prefix, "group", "[GROUP]", "remove", "from", "user", "[USER]")
	h.authToken = authToken
	h.removeGroupToUser = removeGroupToUser
	return h
}

func (h *handler) Match(request *api.Request) bool {
	return h.command.MatchRequest(request) && matcher.MatchRequestAuthToken(h.authToken, request)
}

func (h *handler) Help(request *api.Request) []string {
	if matcher.MatchRequestAuthToken(h.authToken, request) {
		return []string{h.command.Help()}
	}
	return []string{}
}

func (h *handler) HandleMessage(request *api.Request) ([]*api.Response, error) {
	glog.V(2).Infof("handle message: %s", request.Message)
	groupName, err := h.command.Parameter(request, "[GROUP]")
	if err != nil {
		return nil, err
	}
	userName, err := h.command.Parameter(request, "[USER]")
	if err != nil {
		return nil, err
	}
	glog.V(2).Infof("remove group %s to user %s", groupName, userName)
	if err := h.removeGroupToUser(groupName, userName); err != nil {
		glog.V(2).Infof("remove group %s to user %s failed: %v", groupName, userName, err)
		return response.CreateReponseMessage(fmt.Sprintf("remove group %s from user %s failed", groupName, userName)), nil
	}
	glog.V(2).Infof("removed group %s to user %s successful", groupName, userName)
	return response.CreateReponseMessage(fmt.Sprintf("group %s removed from user %s", groupName, userName)), nil
}
