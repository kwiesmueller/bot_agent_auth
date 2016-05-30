package action

import (
	"github.com/bborbe/log"

	"fmt"

	"github.com/bborbe/auth/api"
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

func (a *action) Add(authToken string, token string) error {
	logger.Debugf("add token %s to user with token %s", token, authToken)

	if authToken == token {
		return fmt.Errorf("token equals authToken")
	}

	request := api.AddTokenRequest{
		AuthToken: api.AuthToken(authToken),
		Token:     api.AuthToken(token),
	}
	var response api.AddTokenResponse
	if err := a.callRest("/token", "POST", &request, &response, a.token); err != nil {
		logger.Debugf("add token failed: %v", err)
		return err
	}
	logger.Debugf("add token successful")
	return nil
}
