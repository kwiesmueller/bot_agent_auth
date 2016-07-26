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

func (a *action) Register(authToken string, userName string) error {
	logger.Debugf("register user %s with token %s", userName, authToken)
	request := v1.RegisterRequest{
		AuthToken: model.AuthToken(authToken),
		UserName:  model.UserName(userName),
	}
	var response v1.RegisterResponse
	if err := a.callRest(fmt.Sprintf("/api/1.0/user"), "POST", &request, &response, a.token); err != nil {
		logger.Debugf("register user %s failed: %v", userName, err)
		return err
	}
	logger.Debugf("register user %s successful", userName)
	return nil
}
