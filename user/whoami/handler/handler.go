package handler

import (
	"fmt"
	"strings"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type Whoami func(authToken string) (*api.UserName, error)

type handler struct {
	prefix string
	whoami Whoami
}

func New(prefix string, whoami Whoami) *handler {
	h := new(handler)
	h.prefix = prefix
	h.whoami = whoami
	return h
}

func (h *handler) Match(request *message.Request) bool {
	parts := strings.Split(request.Message, " ")
	return len(parts) == 3 && parts[1] == "application" && parts[2] == "whoami"
}

func (h *handler) Help() string {
	return fmt.Sprintf("%s application whoami [NAME]\n", h.prefix)
}

func (h *handler) HandleMessage(request *message.Request) ([]*message.Response, error) {
	userName, err := h.whoami(request.AuthToken)
	var name string
	if err != nil {
		logger.Debugf("get whoami failed: %v", err)
		name = "-"
	} else {
		name = string(*userName)
	}
	logger.Debugf("application whoamid => send success message")
	return response.CreateReponseMessage(
		fmt.Sprintf("UserToken %s", request.AuthToken),
		fmt.Sprintf("UserName: %s", name),
	), nil
}
