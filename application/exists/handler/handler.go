package handler

import (
	"fmt"
	"strings"

	"github.com/bborbe/bot_agent/message"
	"github.com/bborbe/bot_agent_auth/matcher"
	"github.com/bborbe/bot_agent_auth/response"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type ExistsApplication func(applicationName string) (bool, error)

type handler struct {
	parts             []string
	authToken         string
	existsApplication ExistsApplication
}

func New(prefix string, authToken string, existsApplication ExistsApplication) *handler {
	h := new(handler)
	h.parts = []string{prefix, "application", "exists", "[NAME]"}
	h.authToken = authToken
	h.existsApplication = existsApplication
	return h
}

func (h *handler) Match(request *message.Request) bool {
	return matcher.MatchRequestParts(h.parts, request) && matcher.MatchRequestAuthToken(h.authToken, request)
}

func (h *handler) Help(request *message.Request) []string {
	if matcher.MatchRequestAuthToken(h.authToken, request) {
		return []string{strings.Join(h.parts, " ")}
	}
	return []string{}
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
	if exists {
		return response.CreateReponseMessage(fmt.Sprintf("application %s exists", applicationName)), nil
	} else {
		return response.CreateReponseMessage(fmt.Sprintf("application %s not exists", applicationName)), nil
	}
}
