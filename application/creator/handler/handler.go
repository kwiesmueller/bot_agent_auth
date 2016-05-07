package handler

import (
	"fmt"
	"strings"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent_auth/matcher"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type CreateApplication func(applicationName string) (*api.ApplicationPassword, error)

type handler struct {
	parts             []string
	createApplication CreateApplication
}

func New(prefix string, createApplication CreateApplication) *handler {
	h := new(handler)
	h.parts = []string{prefix, "application", "create", "[NAME]"}
	h.createApplication = createApplication
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
	applicationName := parts[3]
	applicationPassword, err := h.createApplication(applicationName)
	if err != nil {
		logger.Debugf("application creation failed => send failure message: %v", err)
		return response.CreateReponseMessage(fmt.Sprintf("create application %s failed", applicationName)), nil
	}
	logger.Debugf("application created => send success message")
	return response.CreateReponseMessage(fmt.Sprintf("application %s created with password %s", applicationName, *applicationPassword)), nil
}
