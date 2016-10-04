package action

import (
	"github.com/bborbe/bot_agent/api"
	"github.com/golang/glog"

	"fmt"
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

func (a *action) Unregister(authToken api.AuthToken) error {
	glog.V(2).Infof("unregister user with token %s", authToken)
	if err := a.callRest(fmt.Sprintf("/api/1.0/token/%s", authToken), "DELETE", nil, nil, a.token); err != nil {
		glog.V(2).Infof("unregister user with token %s failed: %v", authToken, err)
		return err
	}
	glog.V(2).Infof("unregister user with token %s successful", authToken)
	return nil
}
