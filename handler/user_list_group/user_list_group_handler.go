package user_list_group

import (
	auth_model "github.com/bborbe/auth/model"

	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/golang/glog"
)

type listGroupNamesForUsername func(username auth_model.UserName) ([]auth_model.GroupName, error)

type handler struct {
	command                   command.Command
	listGroupNamesForUsername listGroupNamesForUsername
}

func New(
	prefix model.Prefix,
	lislistGroupNamesForUsernameTokensForUser listGroupNamesForUsername,
) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "user", "[USERNAME]", "list", "groups")
	h.listGroupNamesForUsername = lislistGroupNamesForUsernameTokensForUser
	return h
}

func (h *handler) Allowed(request *api.Request) bool {
	return true
}

func (h *handler) Match(request *api.Request) bool {
	return h.command.MatchRequest(request)
}

func (h *handler) Help(request *api.Request) []string {
	return []string{h.command.Help()}
}

func (h *handler) HandleMessage(request *api.Request) ([]*api.Response, error) {
	glog.V(2).Infof("handle message: %s", request.Message)
	username, err := h.command.Parameter(request, "[USERNAME]")
	if err != nil {
		glog.Warningf("parse parameter failed: %v", err)
		return nil, err
	}
	glog.V(2).Infof("list groups for username %s", username)
	groups, err := h.listGroupNamesForUsername(auth_model.UserName(username))
	if err != nil {
		glog.Warningf("list groups %s failed: %v", username, err)
		return response.CreateReponseMessage("list token failed"), nil
	}
	glog.V(2).Infof("got groups for username %s => send success message", username)
	var results []string
	for _, token := range groups {
		results = append(results, string(token))
	}
	return response.CreateReponseMessage(results...), nil
}
