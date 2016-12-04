package factory

import (
	"net/http"

	"github.com/bborbe/auth/client"
	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/auth/service"
	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/message_handler/allowed_adapter"
	"github.com/bborbe/bot_agent/message_handler/check_tokens_groups"
	"github.com/bborbe/bot_agent/message_handler/match"
	"github.com/bborbe/bot_agent/message_handler/restrict_to_tokens"
	"github.com/bborbe/bot_agent/request_consumer"
	"github.com/bborbe/bot_agent/rest"
	"github.com/bborbe/bot_agent/sender"
	application_creator_handler "github.com/bborbe/bot_agent_auth/handler/application_creator"
	application_deletor_handler "github.com/bborbe/bot_agent_auth/handler/application_deletor"
	application_exists_handler "github.com/bborbe/bot_agent_auth/handler/application_exists"
	token_add_handler "github.com/bborbe/bot_agent_auth/handler/token_add"
	token_remove_handler "github.com/bborbe/bot_agent_auth/handler/token_remove"
	user_create_handler "github.com/bborbe/bot_agent_auth/handler/user_create"
	user_delete_handler "github.com/bborbe/bot_agent_auth/handler/user_delete"
	user_group_add_handler "github.com/bborbe/bot_agent_auth/handler/user_group_add"
	user_group_remove_handler "github.com/bborbe/bot_agent_auth/handler/user_group_remove"
	user_list_handler "github.com/bborbe/bot_agent_auth/handler/user_list"
	user_list_groups_handler "github.com/bborbe/bot_agent_auth/handler/user_list_group"
	user_list_tokens_handler "github.com/bborbe/bot_agent_auth/handler/user_list_token"
	user_register_handler "github.com/bborbe/bot_agent_auth/handler/user_register"
	"github.com/bborbe/bot_agent_auth/handler/user_token_add"
	"github.com/bborbe/bot_agent_auth/handler/user_token_remove"
	user_unregister_handler "github.com/bborbe/bot_agent_auth/handler/user_unregister"
	user_whoami_handler "github.com/bborbe/bot_agent_auth/handler/user_whoami"
	"github.com/bborbe/bot_agent_auth/model"
	http_client_builder "github.com/bborbe/http/client_builder"
	"github.com/bborbe/http/header"
	http_rest "github.com/bborbe/http/rest"
	"github.com/golang/glog"
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

func (b *botAgentAuthfactory) RequestConsumer() request_consumer.RequestConsumer {
	return request_consumer.New(b.sender().Send, b.config.NsqdAddress, b.config.NsqLookupdAddress, b.config.Botname, b.messageHandler())
}

func (b *botAgentAuthfactory) messageHandler() api.MessageHandler {
	var messageHandler api.MessageHandler = match.New(
		b.config.Prefix.String(),
		// without user
		b.userWhoamiHandler(),
		b.userRegisterHandler(),
		// normal user
		b.requireUsername(b.userUnregisterHandler()),
		b.requireUsername(b.tokenAddHandler()),
		b.requireUsername(b.tokenRemoveHandler()),
		// admin user
		b.requrireAdminTokenOrAdminGroup(b.userListTokensHandler()),
		b.requrireAdminTokenOrAdminGroup(b.userListGroupsHandler()),
		b.requrireAdminTokenOrAdminGroup(b.userTokenAddHandler()),
		b.requrireAdminTokenOrAdminGroup(b.userTokenRemoveHandler()),
		b.requrireAdminTokenOrAdminGroup(b.userCreateHandler()),
		b.requrireAdminTokenOrAdminGroup(b.userDeleteHandler()),
		b.requrireAdminTokenOrAdminGroup(b.userAddGroupHandler()),
		b.requrireAdminTokenOrAdminGroup(b.userRemoveGroupHandler()),
		b.requrireAdminTokenOrAdminGroup(b.userListHandler()),
		b.requrireAdminTokenOrAdminGroup(b.applicationCreatorHandler()),
		b.requrireAdminTokenOrAdminGroup(b.applicationDeletorHandler()),
		b.requrireAdminTokenOrAdminGroup(b.applicationExistsHandler()),
	)

	if len(b.config.RestrictToTokens) > 0 {
		messageHandler = restrict_to_tokens.New(
			messageHandler,
			b.config.RestrictToTokens,
		)
	}
	return messageHandler
}

func (b *botAgentAuthfactory) requrireAdminTokenOrAdminGroup(handler match.Handler) match.Handler {
	return check_tokens_groups.New(handler, b.authService().HasGroups, b.config.AdminAuthTokens, b.config.AdminGroups)
}

func (b *botAgentAuthfactory) requireUsername(handler match.Handler) match.Handler {
	return allowed_adapter.New(handler, func(request *api.Request) bool {
		username, err := b.authService().VerifyTokenHasGroups(request.AuthToken, nil)
		if err != nil {
			glog.Warningf("verify token failed: %v", err)
			return false
		}
		return username != nil && len(*username) > 0
	})
}

func (b *botAgentAuthfactory) httpClient() *http.Client {
	return http_client_builder.New().WithoutProxy().Build()
}

func (b *botAgentAuthfactory) authClient() client.Client {
	return client.New(b.httpClient().Do, b.config.AuthUrl, b.config.AuthApplicationName, b.config.AuthApplicationPassword)
}

func (b *botAgentAuthfactory) userService() service.UserService {
	return b.authClient().UserService()
}

func (b *botAgentAuthfactory) authService() service.AuthService {
	return b.authClient().AuthService()
}

func (b *botAgentAuthfactory) userGroupService() service.UserGroupService {
	return b.authClient().UserGroupService()
}

func (b *botAgentAuthfactory) userDataService() service.UserDataService {
	return b.authClient().UserDataService()
}

func (b *botAgentAuthfactory) applicationService() service.ApplicationService {
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

func (b *botAgentAuthfactory) sender() sender.Sender {
	return sender.New(b.producer)
}

func (b *botAgentAuthfactory) whoami(authToken auth_model.AuthToken) (*auth_model.UserName, error) {
	username, err := b.authService().VerifyTokenHasGroups(authToken, nil)
	if err != nil {
		return nil, err
	}
	return username, nil
}

func (b *botAgentAuthfactory) userWhoamiHandler() match.Handler {
	return user_whoami_handler.New(b.config.Prefix, b.whoami)
}

func (b *botAgentAuthfactory) userRegisterHandler() match.Handler {
	return user_register_handler.New(b.config.Prefix, b.userService().CreateUserWithToken)
}

func (b *botAgentAuthfactory) userUnregisterHandler() match.Handler {
	return user_unregister_handler.New(b.config.Prefix, b.userService().DeleteUserWithToken)
}

func (b *botAgentAuthfactory) tokenAddHandler() match.Handler {
	return token_add_handler.New(b.config.Prefix, b.userService().AddTokenToUserWithToken)
}

func (b *botAgentAuthfactory) tokenRemoveHandler() match.Handler {
	return token_remove_handler.New(b.config.Prefix, b.userService().RemoveTokenFromUserWithToken)
}

func (b *botAgentAuthfactory) userListTokensHandler() match.Handler {
	return user_list_tokens_handler.New(b.config.Prefix, b.userService().ListTokenOfUser)
}

func (b *botAgentAuthfactory) userListGroupsHandler() match.Handler {
	return user_list_groups_handler.New(b.config.Prefix, b.userGroupService().ListGroupNamesForUsername)
}

func (b *botAgentAuthfactory) userTokenAddHandler() match.Handler {
	return user_token_add.New(b.config.Prefix, b.userService().AddTokenToUser)
}

func (b *botAgentAuthfactory) userTokenRemoveHandler() match.Handler {
	return user_token_remove.New(b.config.Prefix, b.userService().RemoveTokenFromUser)
}

func (b *botAgentAuthfactory) userCreateHandler() match.Handler {
	return user_create_handler.New(b.config.Prefix, b.userService().CreateUserWithToken)
}

func (b *botAgentAuthfactory) userDeleteHandler() match.Handler {
	return user_delete_handler.New(b.config.Prefix, b.userService().DeleteUser)
}

func (b *botAgentAuthfactory) createApplication(applicationName auth_model.ApplicationName) (*auth_model.ApplicationPassword, error) {
	app, err := b.applicationService().CreateApplication(applicationName)
	if err != nil {
		return nil, err
	}
	return &app.ApplicationPassword, nil
}

func (b *botAgentAuthfactory) applicationCreatorHandler() match.Handler {
	return application_creator_handler.New(b.config.Prefix, b.createApplication)
}

func (b *botAgentAuthfactory) applicationDeletorHandler() match.Handler {
	return application_deletor_handler.New(b.config.Prefix, b.applicationService().DeleteApplication)
}

func (b *botAgentAuthfactory) applicationExistsHandler() match.Handler {
	return application_exists_handler.New(b.config.Prefix, b.applicationService().ExistsApplication)
}

func (b *botAgentAuthfactory) userAddGroupHandler() match.Handler {
	return user_group_add_handler.New(b.config.Prefix, b.userGroupService().AddUserToGroup)
}

func (b *botAgentAuthfactory) userRemoveGroupHandler() match.Handler {
	return user_group_remove_handler.New(b.config.Prefix, b.userGroupService().RemoveUserFromGroup)
}

func (b *botAgentAuthfactory) userListHandler() match.Handler {
	return user_list_handler.New(b.config.Prefix, b.userService().List)
}
