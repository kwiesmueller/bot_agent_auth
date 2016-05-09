package handler

import (
	"fmt"

	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent_auth/command"
	"github.com/bborbe/bot_agent_auth/matcher"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type RemoveGroupToUser func(groupName string, userName string) error

type handler struct {
	command           command.Command
	authToken         string
	removeGroupToUser RemoveGroupToUser
}

func New(prefix string, authToken string, removeGroupToUser RemoveGroupToUser) *handler {
	h := new(handler)
	h.command = command.New(prefix, "remove", "group", "[GROUP]", "from", "user", "[USER]")
	h.authToken = authToken
	h.removeGroupToUser = removeGroupToUser
	return h
}

func (h *handler) Match(request *message.Request) bool {
	return h.command.MatchRequest(request) && matcher.MatchRequestAuthToken(h.authToken, request)
}

func (h *handler) Help(request *message.Request) []string {
	if matcher.MatchRequestAuthToken(h.authToken, request) {
		return []string{h.command.Help()}
	}
	return []string{}
}

func (h *handler) HandleMessage(request *message.Request) ([]*message.Response, error) {
	logger.Debugf("handle message: %s", request.Message)
	groupName, err := h.command.Parameter(request, "[GROUP]")
	if err != nil {
		return nil, err
	}
	userName, err := h.command.Parameter(request, "[USER]")
	if err != nil {
		return nil, err
	}
	logger.Debugf("remove group %s to user %s", groupName, userName)
	if err := h.removeGroupToUser(groupName, userName); err != nil {
		logger.Debugf("remove group %s to user %s failed: %v", groupName, userName, err)
		return response.CreateReponseMessage(fmt.Sprintf("remove group %s from user %s failed", groupName, userName)), nil
	}
	logger.Debugf("removed group %s to user %s successful", groupName, userName)
	return response.CreateReponseMessage(fmt.Sprintf("group %s removed from user %s", groupName, userName)), nil
}
