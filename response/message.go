package response

import "github.com/bborbe/bot_agent/message"

func CreateReponseMessage(msg string) []*message.Response {
	return []*message.Response{&message.Response{
		Message: msg,
		Replay:  false,
	}}
}
