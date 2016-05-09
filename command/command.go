package command

import (
	"github.com/bborbe/bot_agent_auth/matcher"
	"strings"
	"github.com/bborbe/bot_agent/message"
	"fmt"
)

type command struct {
	parts []string
}

type Command interface {
	MatchRequest(request *message.Request) bool
	Parameter(request *message.Request, key string) (string, error)
	Help() string
}

func New(parts... string) *command {
	c := new(command)
	c.parts = parts
	return c
}

func (c *command) Help() string {
	return strings.Join(c.parts, " ")
}

func (c *command) MatchRequest(request *message.Request) bool {
	return matcher.MatchRequestParts(c.parts, request)
}

func (c *command) Parameter(request *message.Request, key string) (string, error) {
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
