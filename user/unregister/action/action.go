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

func (a *action) Unregister(authToken string) error {
	logger.Debugf("unregister user with token %s", authToken)
	request := api.UnRegisterRequest{
		AuthToken: api.AuthToken(authToken),
	}
	var response api.UnRegisterResponse
	if err := a.callRest(fmt.Sprintf("/user"), "DELETE", &request, &response); err != nil {
		logger.Debugf("unregister user with token %s failed: %v", authToken, err)
		return err
	}
	logger.Debugf("unregister user with token %s successful", authToken)
	return nil
}
