package command

import (
	"fmt"
	"strings"

	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent_auth/matcher"
)

type command struct {
	parts []string
}

type Command interface {
	MatchRequest(request *api.Request) bool
	Parameter(request *api.Request, key string) (string, error)
	Help() string
}

func New(parts ...string) *command {
	c := new(command)
	c.parts = parts
	return c
}

func (c *command) Help() string {
	return strings.Join(c.parts, " ")
}

func (c *command) MatchRequest(request *api.Request) bool {
	return matcher.MatchRequestParts(c.parts, request)
}

func (c *command) Parameter(request *api.Request, key string) (string, error) {
	if !c.MatchRequest(request) {
		return "", fmt.Errorf("message does not match command")
	}
	args := strings.Split(request.Message, " ")
	for i, p := range c.parts {
		if p == key {
			return args[i], nil
		}
	}
	return "", fmt.Errorf("parameter %s not found in message", key)
}
