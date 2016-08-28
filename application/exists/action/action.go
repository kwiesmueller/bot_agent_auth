package action

import (
	"fmt"

	"github.com/bborbe/bot_agent/api"

	"github.com/bborbe/auth/v1"
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

func (a *action) Exists(applicationName string) (bool, error) {
	glog.V(2).Infof("exists application %s", applicationName)
	var response v1.GetApplicationResponse
	if err := a.callRest(fmt.Sprintf("/api/1.0/application/%s", applicationName), "GET", nil, &response, a.token); err != nil {
		glog.V(2).Infof("exists application %s failed: %v", applicationName, err)
		return false, err
	}
	glog.V(2).Infof("exists application %s successful", applicationName)
	return len(response.ApplicationPassword) > 0, nil
}
