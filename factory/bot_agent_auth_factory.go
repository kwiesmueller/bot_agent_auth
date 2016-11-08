package factory

import (
	"github.com/bborbe/auth/client/application"
	auth_client_rest "github.com/bborbe/auth/client/rest"
	"github.com/bborbe/auth/client/user"
	"github.com/bborbe/auth/client/user_data"
	"github.com/bborbe/auth/client/user_group"
	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/auth/service"
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/message_handler/match"
	"github.com/bborbe/bot_agent/message_handler/restrict_to_tokens"
	"github.com/bborbe/bot_agent/request_consumer"
	"github.com/bborbe/bot_agent/rest"
	"github.com/bborbe/bot_agent/sender"
	application_creator_handler "github.com/bborbe/bot_agent_auth/application/creator/handler"
	application_deletor_handler "github.com/bborbe/bot_agent_auth/application/deletor/handler"
	application_exists_handler "github.com/bborbe/bot_agent_auth/application/exists/handler"
	"github.com/bborbe/bot_agent_auth/model"
	token_add_handler "github.com/bborbe/bot_agent_auth/token/add/handler"
	token_remove_handler "github.com/bborbe/bot_agent_auth/token/remove/handler"
	user_add_group_handler "github.com/bborbe/bot_agent_auth/user/add_group/handler"
	user_create_handler "github.com/bborbe/bot_agent_auth/user/create/handler"
	user_delete_handler "github.com/bborbe/bot_agent_auth/user/delete/handler"
	user_list_handler "github.com/bborbe/bot_agent_auth/user/list/handler"
	user_register_action "github.com/bborbe/bot_agent_auth/user/register/action"
	user_register_handler "github.com/bborbe/bot_agent_auth/user/register/handler"
	user_remove_group_handler "github.com/bborbe/bot_agent_auth/user/remove_group/handler"
	token_list_handler "github.com/bborbe/bot_agent_auth/user/token_list/handler"
	user_unregister_action "github.com/bborbe/bot_agent_auth/user/unregister/action"
	user_unregister_handler "github.com/bborbe/bot_agent_auth/user/unregister/handler"
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

func (m *botAgentAuthfactory) authClientRest() auth_client_rest.Rest {
	httpClient := http_client_builder.New().WithoutProxy().Build()
	httpRest := http_rest.New(httpClient.Do)
	return auth_client_rest.New(httpRest.Call, m.config.AuthUrl, m.config.AuthApplicationName, m.config.AuthApplicationPassword)
}

func (m *botAgentAuthfactory) UserService() service.UserService {
	return user.New(m.authClientRest().Call)
}

func (m *botAgentAuthfactory) UserGroupService() service.UserGroupService {
	return user_group.New(m.authClientRest().Call)
}

func (m *botAgentAuthfactory) UserDataService() service.UserDataService {
	return user_data.New(m.authClientRest().Call)
}

func (m *botAgentAuthfactory) ApplicationService() service.ApplicationService {
	return application.New(m.authClientRest().Call)
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
	userRegisterAction := user_register_action.New(b.restCaller().Call, b.token())
	return user_register_handler.New(b.config.Prefix, userRegisterAction.Register)
}

func (b *botAgentAuthfactory) userUnregisterHandler() match.Handler {
	userUnregisterAction := user_unregister_action.New(b.restCaller().Call, b.token())
	return user_unregister_handler.New(b.config.Prefix, userUnregisterAction.Unregister)
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

func (b *botAgentAuthfactory) tokenListHandler() match.Handler {
	return token_list_handler.New(b.config.Prefix, b.UserService().ListTokenOfUser)
}
