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
	parts := strings.Split(request.Message, " ")
	return matcher.Match(h.parts, parts)
}

func (h *handler) Help() string {
	return strings.Join(h.parts, " ")
}

func (h *handler) HandleMessage(request *message.Request) ([]*message.Response, error) {
	parts := strings.Split(request.Message, " ")
	userName := parts[3]
	if err := h.register(request.AuthToken, userName); err != nil {
		logger.Debugf("register failed: %v", err)
		return response.CreateReponseMessage("register failed"), nil
	}
	logger.Debugf("user registered => send success message")
	return response.CreateReponseMessage("registration completed"), nil
}
