package action

import (
	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/log"

	"github.com/bborbe/auth/v1"
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

func (a *action) ListUsers() ([]auth_model.UserName, error) {
	logger.Debugf("list users")
	request := v1.UserListRequest{}
	var response v1.UserListResponse
	if err := a.callRest("/api/1.0/user", "GET", &request, &response, a.token); err != nil {
		logger.Debugf("list user failed: %v", err)
		return nil, err
	}
	logger.Debugf("list user successful")

	var result []auth_model.UserName
	for _, user := range response {
		result = append(result, user.UserName)
	}
	return result, nil
}
