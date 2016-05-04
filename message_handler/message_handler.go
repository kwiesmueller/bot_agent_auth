package message_handler

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const PREFIX = "/auth"

type List func(authToken string) ([]string, error)
type Create func(authToken string, authName string) error

type authAgent struct {
	listEntries List
	createEntry Create
}

func New(list List, create Create) *authAgent {
	s := new(authAgent)
	s.listEntries = list
	s.createEntry = create
	return s
}

func (h *authAgent) HandleMessage(request *message.Request) ([]*message.Response, error) {
	logger.Debugf("handle message for token: %v", request.Id)
	if strings.Index(request.Message, PREFIX) != 0 {
		return h.skip()
	}
	parts := strings.Split(request.Message, " ")
	if len(parts) == 2 && parts[1] == "list" {
		return h.list(request.AuthToken)
	}
	if len(parts) == 3 && parts[1] == "create" {
		return h.create(request.AuthToken, parts[2])
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
	fmt.Fprintf(b, "%s list\n", PREFIX)
	fmt.Fprintf(b, "%s create name\n", PREFIX)
	return h.sendMessage(b.String())
}

func (h *authAgent) list(authToken string) ([]*message.Response, error) {
	logger.Debugf("list")
	b := bytes.NewBufferString("")
	list, err := h.listEntries(authToken)
	if err != nil {
		return h.sendMessage("list failed")
	}
	for _, name := range list {
		fmt.Fprintf(b, "%s\n", name)
	}
	return h.sendMessage(b.String())
}

func (h *authAgent) create(authToken string, name string) ([]*message.Response, error) {
	logger.Debugf("create %s", name)
	if err := h.createEntry(authToken, name); err != nil {
		return h.sendMessage(fmt.Sprintf("create %s failed", name))
	}
	return h.sendMessage(fmt.Sprintf("%s created", name))
}

func (h *authAgent) sendMessage(msg string) ([]*message.Response, error) {
	return []*message.Response{&message.Response{
		Message: msg,
		Replay:  false,
	}}, nil
}
