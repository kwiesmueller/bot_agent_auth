package message_handler

import (
	"fmt"
	"strings"

	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/log"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/bot_agent_auth/handler"
	"sort"
)

var logger = log.DefaultLogger

type authAgent struct {
	prefix   string
	handlers []handler.Handler
}

func New(prefix string, handlers... handler.Handler) *authAgent {
	s := new(authAgent)
	s.prefix = prefix
	s.handlers = handlers
	return s
}

func (a *authAgent) HandleMessage(request *message.Request) ([]*message.Response, error) {
	logger.Debugf("handle message for token: %v", request.Id)
	if strings.Index(request.Message, a.prefix) != 0 {
		return a.skip()
	}
	var responses []*message.Response
	for _, h := range a.handlers {
		if h.Match(request) {
			resp, err := h.HandleMessage(request)
			if err != nil {
				return nil, err
			}
			responses = append(responses, resp...)
		}
	}
	if len(responses) == 0 {
		return a.help(), nil
	}
	return responses, nil
}

func (a *authAgent) skip() ([]*message.Response, error) {
	logger.Debugf("message start not with %s => skip", a.prefix)
	return nil, nil
}

func (a *authAgent) help() ([]*message.Response) {
	logger.Debugf("send help message")
	list := []string{fmt.Sprintf("%s help", a.prefix)}
	for _, h := range a.handlers {
		list = append(list, h.Help())
	}
	sort.Strings(list)
	return response.CreateReponseMessage(strings.Join(list, "\n"))
}

