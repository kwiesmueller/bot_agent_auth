package factory

import (
	auth_model "github.com/bborbe/auth/model"

	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/message_handler/match"
	"github.com/bborbe/bot_agent/message_handler/restrict_to_tokens"
	"github.com/bborbe/bot_agent/request_consumer"
	"github.com/bborbe/bot_agent/rest"
	"github.com/bborbe/bot_agent/sender"
	application_creator_action "github.com/bborbe/bot_agent_auth/application/creator/action"
	application_creator_handler "github.com/bborbe/bot_agent_auth/application/creator/handler"
	application_deletor_action "github.com/bborbe/bot_agent_auth/application/deletor/action"
	application_deletor_handler "github.com/bborbe/bot_agent_auth/application/deletor/handler"
	application_exists_action "github.com/bborbe/bot_agent_auth/application/exists/action"
	application_exists_handler "github.com/bborbe/bot_agent_auth/application/exists/handler"
	"github.com/bborbe/bot_agent_auth/model"
	token_add_action "github.com/bborbe/bot_agent_auth/token/add/action"
	token_add_handler "github.com/bborbe/bot_agent_auth/token/add/handler"
	token_remove_action "github.com/bborbe/bot_agent_auth/token/remove/action"
	token_remove_handler "github.com/bborbe/bot_agent_auth/token/remove/handler"
	user_add_group_action "github.com/bborbe/bot_agent_auth/user/add_group/action"
	user_add_group_handler "github.com/bborbe/bot_agent_auth/user/add_group/handler"
	user_create_action "github.com/bborbe/bot_agent_auth/user/create/action"
	user_create_handler "github.com/bborbe/bot_agent_auth/user/create/handler"
	user_delete_action "github.com/bborbe/bot_agent_auth/user/delete/action"
	user_delete_handler "github.com/bborbe/bot_agent_auth/user/delete/handler"
	user_list_action "github.com/bborbe/bot_agent_auth/user/list/action"
	user_list_handler "github.com/bborbe/bot_agent_auth/user/list/handler"
	user_register_action "github.com/bborbe/bot_agent_auth/user/register/action"
	user_register_handler "github.com/bborbe/bot_agent_auth/user/register/handler"
	user_remove_group_action "github.com/bborbe/bot_agent_auth/user/remove_group/action"
	user_remove_group_handler "github.com/bborbe/bot_agent_auth/user/remove_group/handler"
	token_list_action "github.com/bborbe/bot_agent_auth/user/token_list/action"
	token_list_handler "github.com/bborbe/bot_agent_auth/user/token_list/handler"
	user_unregister_action "github.com/bborbe/bot_agent_auth/user/unregister/action"
	user_unregister_handler "github.com/bborbe/bot_agent_auth/user/unregister/handler"
	user_whoami_action "github.com/bborbe/bot_agent_auth/user/whoami/action"
	user_whoami_handler "github.com/bborbe/bot_agent_auth/user/whoami/handler"
	http_client_builder "github.com/bborbe/http/client_builder"
	"github.com/bborbe/http/header"
	http_rest "github.com/bborbe/http/rest"
	"github.com/nsqio/go-nsq"
)

type botAgentAuthfactory struct {
	config   model.Config
	producer *nsq.Producer
}

func New(
	config model.Config,
	producer *nsq.Producer,
) *botAgentAuthfactory {
	b := new(botAgentAuthfactory)
	b.config = config
	b.producer = producer
	return b
}

func (b *botAgentAuthfactory) token() auth_model.AuthToken {
	return auth_model.AuthToken(header.CreateAuthorizationToken(b.config.AuthApplicationName.String(), b.config.AuthApplicationPassword.String()))
}

func (b *botAgentAuthfactory) restCaller() rest.Rest {
	httpClient := http_client_builder.New().WithoutProxy().Build()
	httpRest := http_rest.New(httpClient.Do)
	return rest.New(httpRest.Call, b.config.AuthUrl.String())
}

func (b *botAgentAuthfactory) RequestConsumer() request_consumer.RequestConsumer {

	sender := sender.New(b.producer)

	var messageHandler api.MessageHandler = match.New(
		b.config.Prefix.String(),
		b.applicationCreatorHandler(),
		b.applicationDeletorHandler(),
		b.applicationExistsHandler(),
		b.userWhoamiHandler(),
		b.userRegisterHandler(),
		b.userUnregisterHandler(),
		b.userCreateHandler(),
		b.userDeleteHandler(),
		b.tokenAddHandler(),
		b.tokenRemoveHandler(),
		b.userAddGroupHandler(),
		b.userRemoveGroupHandler(),
		b.userListHandler(),
		b.tokenListHandler(),
	)

	if len(b.config.RestrictToTokens) > 0 {
		messageHandler = restrict_to_tokens.New(
			messageHandler,
			b.config.RestrictToTokens,
		)
	}

	return request_consumer.New(sender.Send, b.config.NsqdAddress, b.config.NsqLookupdAddress, b.config.Botname, messageHandler)
}

func (b *botAgentAuthfactory) applicationCreatorHandler() match.Handler {
	applicationCreatorAction := application_creator_action.New(b.restCaller().Call, b.token())
	return application_creator_handler.New(b.config.Prefix, b.config.AdminAuthToken, applicationCreatorAction.Create)
}
func (b *botAgentAuthfactory) applicationDeletorHandler() match.Handler {
	applicationDeletorAction := application_deletor_action.New(b.restCaller().Call, b.token())
	return application_deletor_handler.New(b.config.Prefix, b.config.AdminAuthToken, applicationDeletorAction.Delete)
}
func (b *botAgentAuthfactory) applicationExistsHandler() match.Handler {
	applicationExistsAction := application_exists_action.New(b.restCaller().Call, b.token())
	return application_exists_handler.New(b.config.Prefix, b.config.AdminAuthToken, applicationExistsAction.Exists)
}
func (b *botAgentAuthfactory) userWhoamiHandler() match.Handler {
	userWhoamiAction := user_whoami_action.New(b.restCaller().Call, b.token())
	return user_whoami_handler.New(b.config.Prefix, userWhoamiAction.Whoami)
}
func (b *botAgentAuthfactory) userRegisterHandler() match.Handler {
	userRegisterAction := user_register_action.New(b.restCaller().Call, b.token())
	return user_register_handler.New(b.config.Prefix, userRegisterAction.Register)
}
func (b *botAgentAuthfactory) userUnregisterHandler() match.Handler {
	userUnregisterAction := user_unregister_action.New(b.restCaller().Call, b.token())
	return user_unregister_handler.New(b.config.Prefix, userUnregisterAction.Unregister)
}
func (b *botAgentAuthfactory) userCreateHandler() match.Handler {
	userCreateAction := user_create_action.New(b.restCaller().Call, b.token())
	return user_create_handler.New(b.config.Prefix, b.config.AdminAuthToken, userCreateAction.CreateUser)
}
func (b *botAgentAuthfactory) userDeleteHandler() match.Handler {
	userDeleteAction := user_delete_action.New(b.restCaller().Call, b.token())
	return user_delete_handler.New(b.config.Prefix, b.config.AdminAuthToken, userDeleteAction.DeleteUser)
}
func (b *botAgentAuthfactory) tokenAddHandler() match.Handler {
	tokenAddAction := token_add_action.New(b.restCaller().Call, b.token())
	return token_add_handler.New(b.config.Prefix, tokenAddAction.Add)
}
func (b *botAgentAuthfactory) tokenRemoveHandler() match.Handler {
	tokenRemoveAction := token_remove_action.New(b.restCaller().Call, b.token())
	return token_remove_handler.New(b.config.Prefix, tokenRemoveAction.Remove)
}
func (b *botAgentAuthfactory) userAddGroupHandler() match.Handler {
	userAddGroupAction := user_add_group_action.New(b.restCaller().Call, b.token())
	return user_add_group_handler.New(b.config.Prefix, b.config.AdminAuthToken, userAddGroupAction.AddGroupToUser)
}
func (b *botAgentAuthfactory) userRemoveGroupHandler() match.Handler {
	userRemoveGroupAction := user_remove_group_action.New(b.restCaller().Call, b.token())
	return user_remove_group_handler.New(b.config.Prefix, b.config.AdminAuthToken, userRemoveGroupAction.RemoveGroupToUser)
}
func (b *botAgentAuthfactory) userListHandler() match.Handler {
	userListAction := user_list_action.New(b.restCaller().Call, b.token())
	return user_list_handler.New(b.config.Prefix, b.config.AdminAuthToken, userListAction.ListUsers)
}
func (b *botAgentAuthfactory) tokenListHandler() match.Handler {
	tokenListAction := token_list_action.New(b.restCaller().Call, b.token())
	return token_list_handler.New(b.config.Prefix, tokenListAction.ListTokensForUser)
}
