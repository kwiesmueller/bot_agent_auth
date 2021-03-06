package user_group_remove

import (
	"fmt"

	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/golang/glog"
)

type removeGroupToUser func(userName auth_model.UserName, groupName auth_model.GroupName) error

type handler struct {
	command           command.Command
	removeGroupToUser removeGroupToUser
}

func New(
	prefix model.Prefix,
	removeGroupToUser removeGroupToUser,
) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "user", "[USERNAME]", "remove", "group", "[GROUP]")
	h.removeGroupToUser = removeGroupToUser
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
	groupName, err := h.command.Parameter(request, "[GROUP]")
	if err != nil {
		return nil, err
	}
	userName, err := h.command.Parameter(request, "[USERNAME]")
	if err != nil {
		return nil, err
	}
	glog.V(2).Infof("remove group %s to user %s", groupName, userName)
	if err := h.removeGroupToUser(auth_model.UserName(userName), auth_model.GroupName(groupName)); err != nil {
		glog.V(2).Infof("remove group %s to user %s failed: %v", groupName, userName, err)
		return response.CreateReponseMessage(fmt.Sprintf("remove group %s from user %s failed", groupName, userName)), nil
	}
	glog.V(2).Infof("removed group %s to user %s successful", groupName, userName)
	return response.CreateReponseMessage(fmt.Sprintf("group %s removed from user %s", groupName, userName)), nil
}
