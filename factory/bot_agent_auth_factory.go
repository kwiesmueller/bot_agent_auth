package factory

import (
	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/auth/service"
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/message_handler/match"
	"github.com/bborbe/bot_agent/message_handler/restrict_to_tokens"
	"github.com/bborbe/bot_agent/request_consumer"
	"github.com/bborbe/bot_agent/rest"
	"github.com/bborbe/bot_agent/sender"
	"github.com/bborbe/bot_agent_auth/model"
	http_client_builder "github.com/bborbe/http/client_builder"
	"github.com/bborbe/http/header"
	http_rest "github.com/bborbe/http/rest"
	"github.com/nsqio/go-nsq"

	"net/http"

	"github.com/bborbe/auth/client"
	application_creator_handler "github.com/bborbe/bot_agent_auth/handler/application_creator"
	application_deletor_handler "github.com/bborbe/bot_agent_auth/handler/application_deletor"
	application_exists_handler "github.com/bborbe/bot_agent_auth/handler/application_exists"
	token_add_handler "github.com/bborbe/bot_agent_auth/handler/token_add"
	token_remove_handler "github.com/bborbe/bot_agent_auth/handler/token_remove"
	user_add_group_handler "github.com/bborbe/bot_agent_auth/handler/user_add_group"
	user_create_handler "github.com/bborbe/bot_agent_auth/handler/user_create"
	user_delete_handler "github.com/bborbe/bot_agent_auth/handler/user_delete"
	user_list_handler "github.com/bborbe/bot_agent_auth/handler/user_list"
	user_list_groups_handler "github.com/bborbe/bot_agent_auth/handler/user_list_group"
	user_list_tokens_handler "github.com/bborbe/bot_agent_auth/handler/user_list_token"
	user_register_handler "github.com/bborbe/bot_agent_auth/handler/user_register"
	user_remove_group_handler "github.com/bborbe/bot_agent_auth/handler/user_remove_group"
	user_unregister_handler "github.com/bborbe/bot_agent_auth/handler/user_unregister"
	user_whoami_handler "github.com/bborbe/bot_agent_auth/handler/user_whoami"
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

func (b *botAgentAuthfactory) httpClient() *http.Client {
	return http_client_builder.New().WithoutProxy().Build()
}

func (b *botAgentAuthfactory) authClient() client.Client {
	return client.New(b.httpClient().Do, b.config.AuthUrl, b.config.AuthApplicationName, b.config.AuthApplicationPassword)
}

func (b *botAgentAuthfactory) UserService() service.UserService {
	return b.authClient().UserService()
}

func (b *botAgentAuthfactory) UserGroupService() service.UserGroupService {
	return b.authClient().UserGroupService()
}

func (b *botAgentAuthfactory) UserDataService() service.UserDataService {
	return b.authClient().UserDataService()
}

func (b *botAgentAuthfactory) ApplicationService() service.ApplicationService {
	return b.authClient().ApplicationService()
}

func (b *botAgentAuthfactory) token() auth_model.AuthToken {
	return auth_model.AuthToken(header.CreateAuthorizationToken(b.config.AuthApplicationName.String(), b.config.AuthApplicationPassword.String()))
}

func (b *botAgentAuthfactory) restCaller() rest.Rest {
	httpClient := http_client_builder.New().WithoutProxy().Build()
	httpRest := http_rest.New(httpClient.Do)
	return rest.New(httpRest.Call, b.config.AuthUrl.String())
}

func (b *botAgentAuthfactory) Sender() sender.Sender {
	return sender.New(b.producer)
}

func (b *botAgentAuthfactory) RequestConsumer() request_consumer.RequestConsumer {
	return request_consumer.New(b.Sender().Send, b.config.NsqdAddress, b.config.NsqLookupdAddress, b.config.Botname, b.MessageHandler())
}

func (b *botAgentAuthfactory) MessageHandler() api.MessageHandler {
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
		b.userListTokensHandler(),
		b.userListGroupsHandler(),
	)

	if len(b.config.RestrictToTokens) > 0 {
		messageHandler = restrict_to_tokens.New(
			messageHandler,
			b.config.RestrictToTokens,
		)
	}
	return messageHandler
}

func (b *botAgentAuthfactory) createApplication(applicationName auth_model.ApplicationName) (*auth_model.ApplicationPassword, error) {
	app, err := b.ApplicationService().CreateApplication(applicationName)
	if err != nil {
		return nil, err
	}
	return &app.ApplicationPassword, nil
}

func (b *botAgentAuthfactory) applicationCreatorHandler() match.Handler {
	return application_creator_handler.New(b.config.Prefix, b.config.AdminAuthToken, b.createApplication)
}

func (b *botAgentAuthfactory) applicationDeletorHandler() match.Handler {
	return application_deletor_handler.New(b.config.Prefix, b.config.AdminAuthToken, b.ApplicationService().DeleteApplication)
}

func (b *botAgentAuthfactory) applicationExistsHandler() match.Handler {
	return application_exists_handler.New(b.config.Prefix, b.config.AdminAuthToken, b.ApplicationService().ExistsApplication)
}

func (b *botAgentAuthfactory) Whoami(authToken auth_model.AuthToken) (*auth_model.UserName, error) {
	username, err := b.UserService().VerifyTokenHasGroups(authToken, nil)
	if err != nil {
		return nil, err
	}
	return username, nil
}

func (b *botAgentAuthfactory) userWhoamiHandler() match.Handler {
	return user_whoami_handler.New(b.config.Prefix, b.Whoami)
}

func (b *botAgentAuthfactory) userRegisterHandler() match.Handler {
	return user_register_handler.New(b.config.Prefix, b.UserService().CreateUserWithToken)
}

func (b *botAgentAuthfactory) userUnregisterHandler() match.Handler {
	return user_unregister_handler.New(b.config.Prefix, b.UserService().DeleteUserWithToken)
}

func (b *botAgentAuthfactory) userCreateHandler() match.Handler {
	return user_create_handler.New(b.config.Prefix, b.config.AdminAuthToken, b.UserService().CreateUserWithToken)
}

func (b *botAgentAuthfactory) userDeleteHandler() match.Handler {
	return user_delete_handler.New(b.config.Prefix, b.config.AdminAuthToken, b.UserService().DeleteUser)
}

func (b *botAgentAuthfactory) tokenAddHandler() match.Handler {
	return token_add_handler.New(b.config.Prefix, b.UserService().AddTokenToUserWithToken)
}

func (b *botAgentAuthfactory) tokenRemoveHandler() match.Handler {
	return token_remove_handler.New(b.config.Prefix, b.UserService().RemoveTokenFromUserWithToken)
}

func (b *botAgentAuthfactory) userAddGroupHandler() match.Handler {
	return user_add_group_handler.New(b.config.Prefix, b.config.AdminAuthToken, b.UserGroupService().AddUserToGroup)
}

func (b *botAgentAuthfactory) userRemoveGroupHandler() match.Handler {
	return user_remove_group_handler.New(b.config.Prefix, b.config.AdminAuthToken, b.UserGroupService().RemoveUserFromGroup)
}

func (b *botAgentAuthfactory) userListHandler() match.Handler {
	return user_list_handler.New(b.config.Prefix, b.config.AdminAuthToken, b.UserService().List)
}

func (b *botAgentAuthfactory) userListTokensHandler() match.Handler {
	return user_list_tokens_handler.New(b.config.Prefix, b.UserService().ListTokenOfUser)
}

func (b *botAgentAuthfactory) userListGroupsHandler() match.Handler {
	return user_list_groups_handler.New(b.config.Prefix, b.UserGroupService().ListGroupNamesForUsername)
}
