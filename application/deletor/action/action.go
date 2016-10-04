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

func (a *action) Delete(applicationName string) error {
	glog.V(2).Infof("delete application %s", applicationName)
	if err := a.callRest(fmt.Sprintf("/api/1.0/application/%s", applicationName), "DELETE", nil, nil, a.token); err != nil {
		glog.V(2).Infof("delete application %s failed: %v", applicationName, err)
		return err
	}
	glog.V(2).Infof("delete application %s successful", applicationName)
	return nil
}
