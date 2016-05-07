package handler

import (
	"fmt"
	"strings"

	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type DeleteApplication func(applicationName string) error

type handler struct {
	prefix            string
	deleteApplication DeleteApplication
}

func New(prefix string, deleteApplication DeleteApplication) *handler {
	h := new(handler)
	h.prefix = prefix
	h.deleteApplication = deleteApplication
	return h
}

func (h *handler) Match(request *message.Request) bool {
	parts := strings.Split(request.Message, " ")
	return len(parts) == 4 && parts[1] == "application" && parts[2] == "delete"
}

func (h *handler) Help() string {
	return fmt.Sprintf("%s application delete [NAME]\n", h.prefix)
}

func (h *handler) HandleMessage(request *message.Request) ([]*message.Response, error) {
	parts := strings.Split(request.Message, " ")
	applicationName := parts[3]
	logger.Debugf("delete applcation %s", applicationName)
	if err := h.deleteApplication(applicationName); err != nil {
		return response.CreateReponseMessage(fmt.Sprintf("delete application %s failed", applicationName)), nil
	}
	return response.CreateReponseMessage(fmt.Sprintf("application %s deleted", applicationName)), nil
}
