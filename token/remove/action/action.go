package action

import (
	"github.com/golang/glog"

	"fmt"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	"github.com/bborbe/bot_agent/api"
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

func (a *action) Remove(authToken api.AuthToken, token api.AuthToken) error {
	glog.V(2).Infof("remove token %s to user with token %s", token, authToken)

	if authToken == token {
		return fmt.Errorf("token equals authToken")
	}

	request := v1.AddTokenRequest{
		AuthToken: model.AuthToken(authToken),
		Token:     model.AuthToken(token),
	}
	var response v1.AddTokenResponse
	if err := a.callRest("/api/1.0/token", "DELETE", &request, &response, a.token); err != nil {
		glog.V(2).Infof("remove token failed: %v", err)
		return err
	}
	glog.V(2).Infof("remove token successful")
	return nil
}
