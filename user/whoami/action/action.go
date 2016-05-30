package action

import (
	"github.com/bborbe/log"

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

func (a *action) Whoami(authToken string) (*api.UserName, error) {
	logger.Debugf("who is %s", authToken)
	request := api.LoginRequest{
		AuthToken: api.AuthToken(authToken),
	}
	var response api.LoginResponse
	if err := a.callRest("/login", "POST", &request, &response, a.token); err != nil {
		logger.Debugf("who is %s failed: %v", authToken, err)
		return nil, err
	}
	logger.Debugf("%s is %v successful", authToken, response.UserName)
	return response.UserName, nil
}
