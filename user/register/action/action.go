package action

import (
	auth_model "github.com/bborbe/auth/model"
	"github.com/golang/glog"

	"fmt"

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

func (a *action) Register(authToken auth_model.AuthToken, userName string) error {
	glog.V(2).Infof("register user %s with token %s", userName, authToken)
	request := v1.RegisterRequest{
		AuthToken: model.AuthToken(authToken),
		UserName:  model.UserName(userName),
	}
	var response v1.RegisterResponse
	if err := a.callRest(fmt.Sprintf("/api/1.0/user"), "POST", &request, &response, a.token); err != nil {
		glog.V(2).Infof("register user %s failed: %v", userName, err)
		return err
	}
	glog.V(2).Infof("register user %s successful", userName)
	return nil
}
