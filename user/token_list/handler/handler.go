package handler

import (
	auth_model "github.com/bborbe/auth/model"

	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/command"
	"github.com/bborbe/bot_agent/response"
	"github.com/bborbe/bot_agent_auth/model"
	"github.com/golang/glog"
)

type listTokensForUser func(username auth_model.UserName) ([]auth_model.AuthToken, error)

type handler struct {
	command           command.Command
	listTokensForUser listTokensForUser
}

func New(prefix model.Prefix, listTokensForUser listTokensForUser) *handler {
	h := new(handler)
	h.command = command.New(prefix.String(), "user", "[USERNAME]", "list", "tokens")
	h.listTokensForUser = listTokensForUser
	return h
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
		return nil, err
	}
	glog.V(2).Infof("list tokens for username %s", username)
	tokens, err := h.listTokensForUser(auth_model.UserName(username))
	if err != nil {
		glog.V(2).Infof("list tokens %s failed: %v", username, err)
		return response.CreateReponseMessage("list token failed"), nil
	}
	glog.V(2).Infof("got tokens for username %s => send success message", username)
	var results []string
	for _, token := range tokens {
		results = append(results, string(token))
	}
	return response.CreateReponseMessage(results...), nil
}
