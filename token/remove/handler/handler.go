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
	h.parts = []string{prefix, "token", "remove", "[NAME]"}
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
	token := parts[3]
	logger.Debugf("remove token %s", token)
	if err := h.register(request.AuthToken, token); err != nil {
		logger.Debugf("remove token %s failed: %v", token, err)
		return response.CreateReponseMessage("remove token failed"), nil
	}
	logger.Debugf("token %s removed => send success message", token)
	return response.CreateReponseMessage("token removed"), nil
}
