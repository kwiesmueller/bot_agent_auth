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

func (a *action) DeleteUser(username string) error {
	logger.Debugf("delete user %s", username)
	if err := a.callRest(fmt.Sprintf("/user/%s", username), "DELETE", nil, nil); err != nil {
		logger.Debugf("delete user %s failed: %v", username, err)
		return err
	}
	logger.Debugf("delete user %s successful", username)
	return nil
}
