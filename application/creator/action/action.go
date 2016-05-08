package action

import (
	"github.com/bborbe/log"

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

func (a *action) Create(applicationName string) (*api.ApplicationPassword, error) {
	logger.Debugf("create application %s", applicationName)
	request := api.CreateApplicationRequest{
		ApplicationName: api.ApplicationName(applicationName),
	}
	var response api.CreateApplicationResponse
	if err := a.callRest("/application", "POST", &request, &response); err != nil {
		logger.Debugf("create application %s failed", applicationName)
		return nil, err
	}
	logger.Debugf("create application %s successful", applicationName)
	return &response.ApplicationPassword, nil
}
