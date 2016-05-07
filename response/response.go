package response

import "github.com/bborbe/bot_agent/message"

func CreateReponseMessage(messages ...string) []*message.Response {
	var result []*message.Response
	for _, msg := range messages {
		result = append(result, &message.Response{
			Message: msg,
			Replay:  false,
		})
	}
	return result
}
