package action

import (
	"github.com/bborbe/log"

	"fmt"

	"github.com/bborbe/auth/api"
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

func (a *action) Remove(authToken string, token string) error {
	logger.Debugf("remove token %s to user with token %s", token, authToken)

	if authToken == token {
		return fmt.Errorf("token equals authToken")
	}

	request := api.AddTokenRequest{
		AuthToken: api.AuthToken(authToken),
		Token:     api.AuthToken(token),
	}
	var response api.AddTokenResponse
	if err := a.callRest("/token", "DELETE", &request, &response); err != nil {
		logger.Debugf("remove token failed: %v", err)
		return err
	}
	logger.Debugf("remove token successful")
	return nil
}
