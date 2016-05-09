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

type DeleteApplication func(applicationName string) error

type handler struct {
	parts             []string
	authToken         string
	deleteApplication DeleteApplication
}

func New(prefix string, authToken string, deleteApplication DeleteApplication) *handler {
	h := new(handler)
	h.parts = []string{prefix, "application", "delete", "[NAME]"}
	h.authToken = authToken
	h.deleteApplication = deleteApplication
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
	logger.Debugf("delete applcation %s", applicationName)
	if err := h.deleteApplication(applicationName); err != nil {
		return response.CreateReponseMessage(fmt.Sprintf("delete application %s failed", applicationName)), nil
	}
	return response.CreateReponseMessage(fmt.Sprintf("application %s deleted", applicationName)), nil
}
