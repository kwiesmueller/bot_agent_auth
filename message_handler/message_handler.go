package message_handler

import (
	"fmt"
	"strings"

	"sort"

	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent_auth/handler"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type authAgent struct {
	prefix   string
	handlers []handler.Handler
}

func New(prefix string, handlers ...handler.Handler) *authAgent {
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
		return a.help(request), nil
	}
	return responses, nil
}

func (a *authAgent) skip() ([]*message.Response, error) {
	logger.Debugf("message start not with %s => skip", a.prefix)
	return nil, nil
}

func (a *authAgent) help(request *message.Request) []*message.Response {
	logger.Debugf("send help message")
	list := []string{fmt.Sprintf("%s help", a.prefix)}
	for _, h := range a.handlers {
		for _, m := range h.Help(request) {
			list = append(list, m)
		}
	}
	sort.Strings(list)
	return response.CreateReponseMessage(strings.Join(list, "\n"))
}
