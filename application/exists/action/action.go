package action

import (
	"github.com/bborbe/log"

	"fmt"

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

func (a *action) Exists(applicationName string) (bool, error) {
	logger.Debugf("exists application %s", applicationName)
	var response api.GetApplicationResponse
	if err := a.callRest(fmt.Sprintf("/application/%s", applicationName), "GET", nil, &response, a.token); err != nil {
		logger.Debugf("exists application %s failed: %v", applicationName, err)
		return false, err
	}
	logger.Debugf("exists application %s successful", applicationName)
	return len(response.ApplicationPassword) > 0, nil
}
