package handler

import (
	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent_auth/command"
	"github.com/bborbe/bot_agent_auth/matcher"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type DeleteUser func(username string) error

type handler struct {
	command   command.Command
	authToken string
	delete    DeleteUser
}

func New(prefix string, authToken string, delete DeleteUser) *handler {
	h := new(handler)
	h.command = command.New(prefix, "user", "delete", "[USERNAME]")
	h.authToken = authToken
	h.delete = delete
	return h
}

func (h *handler) Match(request *message.Request) bool {
	return h.command.MatchRequest(request) && matcher.MatchRequestAuthToken(h.authToken, request)
}

func (h *handler) Help(request *message.Request) []string {
	return []string{h.command.Help()}
}

func (h *handler) HandleMessage(request *message.Request) ([]*message.Response, error) {
	logger.Debugf("delete user")
	username, err := h.command.Parameter(request, "[USERNAME]")
	if err != nil {
		return nil, err
	}
	logger.Debugf("delete user %s", username)
	if err := h.delete(username); err != nil {
		logger.Debugf("delete user %s failed: %v", username, err)
		return response.CreateReponseMessage("delete failed"), nil
	}
	logger.Debugf("delete user %s successful", username)
	return response.CreateReponseMessage("delete user completed"), nil
}
