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

func (a *action) Create(applicationName string) (*api.ApplicationPassword, error) {
	logger.Debugf("create application %s", applicationName)
	request := api.CreateApplicationRequest{
		ApplicationName: api.ApplicationName(applicationName),
	}
	var response api.CreateApplicationResponse
	if err := a.callRest("/application", "POST", &request, &response, a.token); err != nil {
		logger.Debugf("create application %s failed: %v", applicationName, err)
		return nil, err
	}
	logger.Debugf("create application %s successful", applicationName)
	return &response.ApplicationPassword, nil
}
