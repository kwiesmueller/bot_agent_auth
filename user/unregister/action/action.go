package action

import (
	"github.com/bborbe/log"

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

func (a *action) Unregister(authToken string) error {
	logger.Debugf("unregister user with token %s", authToken)
	if err := a.callRest(fmt.Sprintf("/user/%s", authToken), "DELETE", nil, nil); err != nil {
		logger.Debugf("unregister user with token %s failed: %v", authToken, err)
		return err
	}
	logger.Debugf("unregister user with token %s successful", authToken)
	return nil
}
