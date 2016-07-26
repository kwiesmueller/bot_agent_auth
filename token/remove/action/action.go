package action

import (
	"github.com/bborbe/log"

	"fmt"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
)

var logger = log.DefaultLogger

type CallRest func(path string, method string, request interface{}, response interface{}, token string) error

type action struct {
	callRest CallRest
	token    string
}

func New(callRest CallRest, token string) *action {
	m := new(action)
	m.callRest = callRest
	m.token = token
	return m
}

func (a *action) Remove(authToken string, token string) error {
	logger.Debugf("remove token %s to user with token %s", token, authToken)

	if authToken == token {
		return fmt.Errorf("token equals authToken")
	}

	request := v1.AddTokenRequest{
		AuthToken: model.AuthToken(authToken),
		Token:     model.AuthToken(token),
	}
	var response v1.AddTokenResponse
	if err := a.callRest("/api/v1.0/token", "DELETE", &request, &response, a.token); err != nil {
		logger.Debugf("remove token failed: %v", err)
		return err
	}
	logger.Debugf("remove token successful")
	return nil
}
