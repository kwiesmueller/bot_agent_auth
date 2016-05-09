package action

import (
	"github.com/bborbe/log"

	"github.com/bborbe/auth/api"
)

var logger = log.DefaultLogger

type CallRest func(path string, method string, request interface{}, response interface{}) error

type action struct {
	callRest CallRest
}

func New(callRest CallRest) *action {
	m := new(action)
	m.callRest = callRest
	return m
}

func (a *action) RemoveGroupToUser(groupName string, userName string) error {
	logger.Debugf("remove user %s from group %s", userName, groupName)
	request := api.AddUserToGroupRequest{
		UserName:  api.UserName(userName),
		GroupName: api.GroupName(groupName),
	}
	var response api.AddUserToGroupResponse
	if err := a.callRest("/user_group", "DELETE", &request, &response); err != nil {
		logger.Debugf("remove user %v from group %v failed: %v", userName, groupName, err)
		return err
	}
	logger.Debugf("remove user user %v from group %v successful", userName, groupName)
	return nil
}
