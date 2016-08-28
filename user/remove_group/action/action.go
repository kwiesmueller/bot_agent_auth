package action

import (
	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	"github.com/bborbe/bot_agent/api"
	"github.com/golang/glog"
)

type CallRest func(path string, method string, request interface{}, response interface{}, token api.AuthToken) error

type action struct {
	callRest CallRest
	token    api.AuthToken
}

func New(callRest CallRest, token api.AuthToken) *action {
	m := new(action)
	m.callRest = callRest
	m.token = token
	return m
}

func (a *action) RemoveGroupToUser(groupName string, userName string) error {
	glog.V(2).Infof("remove user %s from group %s", userName, groupName)
	request := v1.AddUserToGroupRequest{
		UserName:  model.UserName(userName),
		GroupName: model.GroupName(groupName),
	}
	var response v1.AddUserToGroupResponse
	if err := a.callRest("/api/1.0/user_group", "DELETE", &request, &response, a.token); err != nil {
		glog.V(2).Infof("remove user %v from group %v failed: %v", userName, groupName, err)
		return err
	}
	glog.V(2).Infof("remove user user %v from group %v successful", userName, groupName)
	return nil
}
