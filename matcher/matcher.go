package matcher

import (
	"strings"

	"github.com/bborbe/bot_agent/message"
)

func MatchRequestParts(requiredParts []string, request *message.Request) bool {
	parts := strings.Split(request.Message, " ")
	return MatchParts(requiredParts, parts)
}

func MatchParts(requiredParts []string, parts []string) bool {
	if len(requiredParts) != len(parts) {
		return false
	}
	for i, _ := range requiredParts {

		if requiredParts[i] != parts[i] {
			if len(requiredParts[i]) <= 2 ||
				requiredParts[i][0:1] != "[" ||
				requiredParts[i][len(requiredParts[i])-1:len(requiredParts[i])] != "]" {
				return false
			}
		}
	}
	return true
}

func MatchRequestAuthToken(requiredAuthToken string, request *message.Request) bool {
	return MatchAuthToken(requiredAuthToken, request.AuthToken)
}

func MatchAuthToken(requiredAuthToken string, authToken string) bool {
	return requiredAuthToken == authToken
}
