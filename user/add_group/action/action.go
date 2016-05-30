package action

import (
	"github.com/bborbe/log"

	"github.com/bborbe/auth/api"
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

func (a *action) AddGroupToUser(groupName string, userName string) error {
	logger.Debugf("add user %s to group %s", userName, groupName)
	request := api.AddUserToGroupRequest{
		UserName:  api.UserName(userName),
		GroupName: api.GroupName(groupName),
	}
	var response api.AddUserToGroupResponse
	if err := a.callRest("/user_group", "POST", &request, &response, a.token); err != nil {
		logger.Debugf("add user %v to group %v failed: %v", userName, groupName, err)
		return err
	}
	logger.Debugf("add user user %v to group %v successful", userName, groupName)
	return nil
}
