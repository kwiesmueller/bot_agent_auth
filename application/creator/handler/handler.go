package handler

import (
	"fmt"
	"github.com/bborbe/log"
	"github.com/bborbe/bot_agent/message"
	"strings"
	"github.com/bborbe/auth/api"
	"github.com/bborbe/bot_agent_auth/response"
)

var logger = log.DefaultLogger

type CreateApplication func(applicationName string) (*api.ApplicationPassword, error)

type handler struct {
	prefix            string
	createApplication CreateApplication
}

func New(prefix string, createApplication CreateApplication) *handler {
	h := new(handler)
	h.prefix = prefix
	h.createApplication = createApplication
	return h
}

func (h *handler) Match(request *message.Request) bool {
	parts := strings.Split(request.Message, " ")
	return len(parts) == 4 && parts[1] == "application" && parts[2] == "create"
}

func (h *handler) Help() string {
	return fmt.Sprintf("%s application create [NAME]\n", h.prefix)
}

func (h *handler) HandleMessage(request *message.Request) ([]*message.Response, error) {
	parts := strings.Split(request.Message, " ")
	applicationName := parts[3]
	applicationPassword, err := h.createApplication(applicationName)
	if err != nil {
		logger.Debugf("application creation failed => send failure message: %v", err)
		return response.CreateReponseMessage(fmt.Sprintf("create application %s failed", applicationName)),nil
	}
	logger.Debugf("application created => send success message")
	return response.CreateReponseMessage(fmt.Sprintf("application %s created with password %s", applicationName, *applicationPassword)),nil
}
