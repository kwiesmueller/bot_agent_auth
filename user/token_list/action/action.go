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

func (a *action) ListTokensForUser(username auth_model.UserName) ([]auth_model.AuthToken, error) {
	glog.V(2).Infof("list tokens for username %s", username)
	response := []auth_model.AuthToken{}
	if err := a.callRest(fmt.Sprintf("/api/1.0/token?username=%v", username), "GET", nil, &response, a.token); err != nil {
		glog.V(2).Infof("list token failed: %v", err)
		return nil, err
	}
	glog.V(2).Infof("list token successful")
	return response, nil
}
