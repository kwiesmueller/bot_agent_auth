package handler

import (
	"fmt"
	"github.com/bborbe/log"
	"github.com/bborbe/bot_agent/message"
	"strings"
	"github.com/bborbe/bot_agent_auth/response"
)

var logger = log.DefaultLogger

type ExistsApplication func(applicationName string) (*bool, error)

type handler struct {
	prefix            string
	existsApplication ExistsApplication
}

func New(prefix string, existsApplication ExistsApplication) *handler {
	h := new(handler)
	h.prefix = prefix
	h.existsApplication = existsApplication
	return h
}

func (h *handler) Match(request *message.Request) bool {
	parts := strings.Split(request.Message, " ")
	return len(parts) == 4 && parts[1] == "application" && parts[2] == "exists"
}

func (h *handler) Help() string {
	return fmt.Sprintf("%s application exists [NAME]\n", h.prefix)
}

func (h *handler) HandleMessage(request *message.Request) ([]*message.Response, error) {
	parts := strings.Split(request.Message, " ")
	applicationName := parts[3]
	exists, err := h.existsApplication(applicationName)
	if err != nil {
		logger.Debugf("application exists failed => send failure message: %v", err)
		return response.CreateReponseMessage(fmt.Sprintf("exists application %s failed", applicationName)), nil
	}
	logger.Debugf("application exists => send success message")
	if *exists {
		return response.CreateReponseMessage(fmt.Sprintf("application %s exists", applicationName)), nil
	} else {
		return response.CreateReponseMessage(fmt.Sprintf("application %s not exists", applicationName)), nil
	}
}
