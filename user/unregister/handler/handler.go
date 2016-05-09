package handler

import (
	"strings"

	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent_auth/matcher"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type Unregister func(authToken string) error

type handler struct {
	parts      []string
	unregister Unregister
}

func New(prefix string, unregister Unregister) *handler {
	h := new(handler)
	h.parts = []string{prefix, "unregister"}
	h.unregister = unregister
	return h
}

func (h *handler) Match(request *message.Request) bool {
	return matcher.MatchRequestParts(h.parts, request)
}

func (h *handler) Help(request *message.Request) []string {
	return []string{strings.Join(h.parts, " ")}
}

func (h *handler) HandleMessage(request *message.Request) ([]*message.Response, error) {
	logger.Debugf("unregister user with token %s", request.AuthToken)
	if err := h.unregister(request.AuthToken); err != nil {
		logger.Debugf("unregister user with token %s failed: %v", request.AuthToken, err)
		return response.CreateReponseMessage("unregister failed"), nil
	}
	logger.Debugf("unregister user with token %s successful", request.AuthToken)
	return response.CreateReponseMessage("unregistration completed"), nil
}
