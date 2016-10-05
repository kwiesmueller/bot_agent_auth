package action

import (
	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	"github.com/golang/glog"
)

type CallRest func(path string, method string, request interface{}, response interface{}, token auth_model.AuthToken) error

type action struct {
	callRest CallRest
	token    auth_model.AuthToken
}

func New(callRest CallRest, token auth_model.AuthToken) *action {
	m := new(action)
	m.callRest = callRest
	m.token = token
	return m
}

func (a *action) ListUsers() ([]auth_model.UserName, error) {
	glog.V(2).Infof("list users")
	request := v1.UserListRequest{}
	var response v1.UserListResponse
	if err := a.callRest("/api/1.0/user", "GET", &request, &response, a.token); err != nil {
		glog.V(2).Infof("list user failed: %v", err)
		return nil, err
	}
	glog.V(2).Infof("list user successful")

	var result []auth_model.UserName
	for _, user := range response {
		result = append(result, user.UserName)
	}
	return result, nil
}
