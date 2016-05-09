package handler

import (
	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
	"github.com/bborbe/bot_agent_auth/command"
	"github.com/bborbe/bot_agent_auth/matcher"
	"fmt"
)

var logger = log.DefaultLogger

type AddGroupToUser func(groupName string, userName string) error

type handler struct {
	command        command.Command
	authToken      string
	addGroupToUser AddGroupToUser
}

func New(prefix string, authToken string, addGroupToUser AddGroupToUser) *handler {
	h := new(handler)
	h.command = command.New(prefix, "add", "group", "[GROUP]", "to", "user", "[USER]")
	h.authToken = authToken
	h.addGroupToUser = addGroupToUser
	return h
}

func (h *handler) Match(request *message.Request) bool {
	return h.command.MatchRequest(request)   && matcher.MatchRequestAuthToken(h.authToken, request)
}

func (h *handler) Help(request *message.Request) []string {
	return []string{h.command.Help()}
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
	logger.Debugf("add group %s to user %s", groupName, userName)
	if err := h.addGroupToUser(groupName, userName); err != nil {
		logger.Debugf("add group %s to user %s failed: %v", groupName, userName, err)
		return response.CreateReponseMessage(fmt.Sprintf("add group %s to user %s failed", groupName, userName)), nil
	}
	logger.Debugf("added group %s to user %s successful", groupName, userName)
	return response.CreateReponseMessage(fmt.Sprintf("group %s added to %s", groupName, userName)), nil
}
