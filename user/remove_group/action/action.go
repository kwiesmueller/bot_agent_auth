package action

import (
	"github.com/bborbe/log"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
)

var logger = log.DefaultLogger

type CallRest func(path string, method string, request interface{}, response interface{}, token string) error

type action struct {
	callRest CallRest
	token    string
}

func New(callRest CallRest, token string) *action {
	m := new(action)
	m.callRest = callRest
	m.token = token
	return m
}

func (a *action) RemoveGroupToUser(groupName string, userName string) error {
	logger.Debugf("remove user %s from group %s", userName, groupName)
	request := v1.AddUserToGroupRequest{
		UserName:  model.UserName(userName),
		GroupName: model.GroupName(groupName),
	}
	var response v1.AddUserToGroupResponse
	if err := a.callRest("/api/v1.0/user_group", "DELETE", &request, &response, a.token); err != nil {
		logger.Debugf("remove user %v from group %v failed: %v", userName, groupName, err)
		return err
	}
	logger.Debugf("remove user user %v from group %v successful", userName, groupName)
	return nil
}
