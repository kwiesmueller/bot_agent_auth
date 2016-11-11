package user_group_add

import (
	"fmt"

	auth_model "github.com/bborbe/auth/model"

	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/matcher"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/golang/glog"
)

type AddGroupToUser func(userName auth_model.UserName, groupName auth_model.GroupName) error

type handler struct {
	command        command.Command
	authToken      auth_model.AuthToken
	addGroupToUser AddGroupToUser
}

func New(prefix model.Prefix, authToken auth_model.AuthToken, addGroupToUser AddGroupToUser) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "user", "[USERNAME]", "add", "group", "[GROUP]")
	h.authToken = authToken
	h.addGroupToUser = addGroupToUser
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
	userName, err := h.command.Parameter(request, "[USERNAME]")
	if err != nil {
		return nil, err
	}
	glog.V(2).Infof("add group %s to user %s", groupName, userName)
	if err := h.addGroupToUser(auth_model.UserName(userName), auth_model.GroupName(groupName)); err != nil {
		glog.V(2).Infof("add group %s to user %s failed: %v", groupName, userName, err)
		return response.CreateReponseMessage(fmt.Sprintf("add group %s to user %s failed", groupName, userName)), nil
	}
	glog.V(2).Infof("added group %s to user %s successful", groupName, userName)
	return response.CreateReponseMessage(fmt.Sprintf("group %s added to %s", groupName, userName)), nil
}
