package message_handler

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const PREFIX = "/auth"

type CreateApplication func(applicationName string) (*api.ApplicationPassword, error)
type DeleteApplication func(applicationName string) error
type ExistsApplication func(applicationName string) (*bool, error)

type authAgent struct {
	createApplication CreateApplication
	deleteApplication DeleteApplication
	existsApplication ExistsApplication
}

func New(createApplication CreateApplication, deleteApplication DeleteApplication, existsApplication ExistsApplication) *authAgent {
	s := new(authAgent)
	s.createApplication = createApplication
	s.deleteApplication = deleteApplication
	s.existsApplication = existsApplication
	return s
}

func (h *authAgent) HandleMessage(request *message.Request) ([]*message.Response, error) {
	logger.Debugf("handle message for token: %v", request.Id)
	if strings.Index(request.Message, PREFIX) != 0 {
		return h.skip()
	}
	parts := strings.Split(request.Message, " ")
	if len(parts) == 4 && parts[1] == "application" && parts[2] == "create" {
		applicationName := parts[3]
		applicationPassword, err := h.createApplication(applicationName)
		if err != nil {
			logger.Debugf("application creation failed => send failure message: %v", err)
			return h.sendMessage(fmt.Sprintf("create application %s failed", applicationName))
		}
		logger.Debugf("application created => send success message")
		return h.sendMessage(fmt.Sprintf("application %s created with password %s", applicationName, *applicationPassword))
	}
	if len(parts) == 4 && parts[1] == "application" && parts[2] == "delete" {
		applicationName := parts[3]
		logger.Debugf("delete applcation %s", applicationName)
		if err := h.deleteApplication(applicationName); err != nil {
			return h.sendMessage(fmt.Sprintf("delete application %s failed", applicationName))
		}
		return h.sendMessage(fmt.Sprintf("application %s deleted", applicationName))
	}
	if len(parts) == 4 && parts[1] == "application" && parts[2] == "exists" {
		applicationName := parts[3]
		exists, err := h.existsApplication(applicationName)
		if err != nil {
			logger.Debugf("application exists failed => send failure message: %v", err)
			return h.sendMessage(fmt.Sprintf("create application %s failed", applicationName))
		}
		logger.Debugf("application exists => send success message")
		if *exists {
			return h.sendMessage(fmt.Sprintf("application %s exists", applicationName))
		} else {
			return h.sendMessage(fmt.Sprintf("application %s not exists", applicationName))
		}
	}
	return h.help()
}

func (h *authAgent) skip() ([]*message.Response, error) {
	logger.Debugf("message start not with %s => skip", PREFIX)
	return nil, nil
}

func (h *authAgent) help() ([]*message.Response, error) {
	logger.Debugf("send help message")
	b := bytes.NewBufferString("")
	fmt.Fprintf(b, "%s help\n", PREFIX)
	fmt.Fprintf(b, "%s application create [NAME]\n", PREFIX)
	fmt.Fprintf(b, "%s application delete [NAME]\n", PREFIX)
	fmt.Fprintf(b, "%s application exists [NAME]\n", PREFIX)
	return h.sendMessage(b.String())
}

func (h *authAgent) sendMessage(msg string) ([]*message.Response, error) {
	return []*message.Response{&message.Response{
		Message: msg,
		Replay:  false,
	}}, nil
}
