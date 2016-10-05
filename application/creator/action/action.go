package action

import (
	auth_model "github.com/bborbe/auth/model"
	"github.com/golang/glog"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
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

func (a *action) Create(applicationName string) (*model.ApplicationPassword, error) {
	glog.V(2).Infof("create application %s", applicationName)
	request := v1.CreateApplicationRequest{
		ApplicationName: model.ApplicationName(applicationName),
	}
	var response v1.CreateApplicationResponse
	if err := a.callRest("/api/1.0/application", "POST", &request, &response, a.token); err != nil {
		glog.V(2).Infof("create application %s failed: %v", applicationName, err)
		return nil, err
	}
	glog.V(2).Infof("create application %s successful", applicationName)
	return &response.ApplicationPassword, nil
}
