package handler

import (
	"strings"

	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent_auth/matcher"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type Register func(authToken string, userName string) error

type handler struct {
	parts    []string
	register Register
}

func New(prefix string, register Register) *handler {
	h := new(handler)
	h.parts = []string{prefix, "register", "[NAME]"}
	h.register = register
	return h
}

func (h *handler) Match(request *message.Request) bool {
	return matcher.MatchRequestParts(h.parts, request)
}

func (h *handler) Help(request *message.Request) []string {
	return []string{strings.Join(h.parts, " ")}
}

func (h *handler) HandleMessage(request *message.Request) ([]*message.Response, error) {
	parts := strings.Split(request.Message, " ")
	userName := parts[2]
	logger.Debugf("register user %s", userName)
	if err := h.register(request.AuthToken, userName); err != nil {
		logger.Debugf("register %s failed: %v", userName, err)
		return response.CreateReponseMessage("register failed"), nil
	}
	logger.Debugf("user %s  registered => send success message", userName)
	return response.CreateReponseMessage("registration completed"), nil
}
