package action

import (
	"github.com/bborbe/log"

	"github.com/bborbe/auth/api"
	"fmt"
)

var logger = log.DefaultLogger

type CallRest func(path string, method string, request interface{}, response interface{}) error

type action struct {
	callRest CallRest
}

func New(callRest CallRest) *action {
	m := new(action)
	m.callRest = callRest
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
	if err := a.callRest("/token", "POST", &request, &response); err != nil {
		logger.Debugf("add token failed: %v", err)
		return err
	}
	logger.Debugf("add token successful")
	return nil
}
