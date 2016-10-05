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

func (a *action) Whoami(authToken auth_model.AuthToken) (*model.UserName, error) {
	glog.V(2).Infof("who is %s", authToken)
	request := v1.LoginRequest{
		AuthToken: model.AuthToken(authToken),
	}
	var response v1.LoginResponse
	if err := a.callRest("/api/1.0/login", "POST", &request, &response, a.token); err != nil {
		glog.V(2).Infof("who is %s failed: %v", authToken, err)
		return nil, err
	}
	glog.V(2).Infof("%s is %v successful", authToken, response.UserName)
	return response.UserName, nil
}
