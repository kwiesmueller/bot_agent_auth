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

func (a *action) Delete(applicationName string) error {
	logger.Debugf("delete application %s", applicationName)
	if err := a.callRest(fmt.Sprintf("/application/%s", applicationName), "DELETE", nil, nil, a.token); err != nil {
		logger.Debugf("delete application %s failed: %v", applicationName, err)
		return err
	}
	logger.Debugf("delete application %s successful", applicationName)
	return nil
}
