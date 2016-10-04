package action

import (
	"fmt"

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

func (a *action) DeleteUser(username string) error {
	glog.V(2).Infof("delete user %s", username)
	if err := a.callRest(fmt.Sprintf("/api/1.0/user/%s", username), "DELETE", nil, nil, a.token); err != nil {
		glog.V(2).Infof("delete user %s failed: %v", username, err)
		return err
	}
	glog.V(2).Infof("delete user %s successful", username)
	return nil
}
