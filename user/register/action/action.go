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

func (a *action) Register(authToken string, userName string) error {
	logger.Debugf("register user %s with token %s", userName, authToken)
	request := api.RegisterRequest{
		AuthToken: api.AuthToken(authToken),
		UserName:  api.UserName(userName),
	}
	var response api.RegisterResponse
	if err := a.callRest(fmt.Sprintf("/user"), "POST", &request, &response); err != nil {
		logger.Debugf("register user %s failed: %v", userName, err)
		return err
	}
	logger.Debugf("register user %s successful", userName)
	return nil
}
