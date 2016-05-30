package action

import (
	"github.com/bborbe/log"

	"fmt"
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

func (a *action) Unregister(authToken string) error {
	logger.Debugf("unregister user with token %s", authToken)
	if err := a.callRest(fmt.Sprintf("/token/%s", authToken), "DELETE", nil, nil, a.token); err != nil {
		logger.Debugf("unregister user with token %s failed: %v", authToken, err)
		return err
	}
	logger.Debugf("unregister user with token %s successful", authToken)
	return nil
}
