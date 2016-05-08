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

func (a *action) Delete(applicationName string) error {
	logger.Debugf("delete application %s", applicationName)
	var request api.DeleteApplicationRequest
	var response api.DeleteApplicationResponse
	if err := a.callRest(fmt.Sprintf("/application/%s", applicationName), "DELETE", &request, &response); err != nil {
		logger.Debugf("delete application %s failed: %v", applicationName, err)
		return err
	}
	logger.Debugf("delete application %s successful", applicationName)
	return nil
}
