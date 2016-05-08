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

func (a *action) Exists(applicationName string) (bool, error) {
	logger.Debugf("exists application %s", applicationName)
	var request api.GetApplicationRequest
	var response api.GetApplicationResponse
	if err := a.callRest(fmt.Sprintf("/application/%s", applicationName), "GET", &request, &response); err != nil {
		logger.Debugf("exists application %s failed: %v", applicationName, err)
		return false, err
	}
	logger.Debugf("exists application %s successful", applicationName)
	return len(response.ApplicationPassword) > 0, nil
}
