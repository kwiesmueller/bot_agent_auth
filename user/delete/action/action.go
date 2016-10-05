package action

import (
	"fmt"
	auth_model "github.com/bborbe/auth/model"

	"github.com/golang/glog"
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

func (a *action) DeleteUser(username string) error {
	glog.V(2).Infof("delete user %s", username)
	if err := a.callRest(fmt.Sprintf("/api/1.0/user/%s", username), "DELETE", nil, nil, a.token); err != nil {
		glog.V(2).Infof("delete user %s failed: %v", username, err)
		return err
	}
	glog.V(2).Infof("delete user %s successful", username)
	return nil
}
