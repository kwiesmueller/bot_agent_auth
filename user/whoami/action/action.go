package action

import (
	"github.com/bborbe/log"

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

func (a *action) Whoami(authToken string) (*model.UserName, error) {
	logger.Debugf("who is %s", authToken)
	request := v1.LoginRequest{
		AuthToken: model.AuthToken(authToken),
	}
	var response v1.LoginResponse
	if err := a.callRest("/api/v1.0/api/v1.0/login", "POST", &request, &response, a.token); err != nil {
		logger.Debugf("who is %s failed: %v", authToken, err)
		return nil, err
	}
	logger.Debugf("%s is %v successful", authToken, response.UserName)
	return response.UserName, nil
}
