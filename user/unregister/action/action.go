package action

import (
	auth_model "github.com/bborbe/auth/model"
	"github.com/golang/glog"

	"fmt"
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

func (a *action) Unregister(authToken auth_model.AuthToken) error {
	glog.V(2).Infof("unregister user with token %s", authToken)
	if err := a.callRest(fmt.Sprintf("/api/1.0/token/%s", authToken), "DELETE", nil, nil, a.token); err != nil {
		glog.V(2).Infof("unregister user with token %s failed: %v", authToken, err)
		return err
	}
	glog.V(2).Infof("unregister user with token %s successful", authToken)
	return nil
}
