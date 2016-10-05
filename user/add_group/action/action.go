package action

import (
	"github.com/bborbe/auth/model"
	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	"github.com/golang/glog"
)

type CallRest func(path string, method string, request interface{}, response interface{}, token auth_model.AuthToken) error

type action struct {
	callRest CallRest
	token    auth_model.AuthToken
}

func New(callRest CallRest, token auth_model.AuthToken) *action {
	m := new(action)
	m.callRest = callRest
	m.token = token
	return m
}

func (a *action) AddGroupToUser(groupName string, userName string) error {
	glog.V(2).Infof("add user %s to group %s", userName, groupName)
	request := v1.AddUserToGroupRequest{
		UserName:  model.UserName(userName),
		GroupName: model.GroupName(groupName),
	}
	var response v1.AddUserToGroupResponse
	if err := a.callRest("/api/1.0/user_group", "POST", &request, &response, a.token); err != nil {
		glog.V(2).Infof("add user %v to group %v failed: %v", userName, groupName, err)
		return err
	}
	glog.V(2).Infof("add user user %v to group %v successful", userName, groupName)
	return nil
}
