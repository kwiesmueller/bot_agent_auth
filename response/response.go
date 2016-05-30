package response

import "github.com/bborbe/bot_agent/api"

func CreateReponseMessage(messages ...string) []*api.Response {
	var result []*api.Response
	for _, msg := range messages {
		result = append(result, &api.Response{
			Message: msg,
			Replay:  false,
		})
	}
	return result
}
