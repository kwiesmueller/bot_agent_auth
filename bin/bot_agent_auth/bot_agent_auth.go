package main

import (
	"fmt"
	"os"

	"strings"

	"github.com/bborbe/bot_agent/api"
	"github.com/bborbe/bot_agent/message_handler/match"
	"github.com/bborbe/bot_agent/message_handler/restrict_to_tokens"
	"github.com/bborbe/bot_agent/producer"
	"github.com/bborbe/bot_agent/request_consumer"
	"github.com/bborbe/bot_agent/rest"
	"github.com/bborbe/bot_agent/sender"
	application_creator_action "github.com/bborbe/bot_agent_auth/application/creator/action"
	application_creator_handler "github.com/bborbe/bot_agent_auth/application/creator/handler"
	application_deletor_action "github.com/bborbe/bot_agent_auth/application/deletor/action"
	application_deletor_handler "github.com/bborbe/bot_agent_auth/application/deletor/handler"
	application_exists_action "github.com/bborbe/bot_agent_auth/application/exists/action"
	application_exists_handler "github.com/bborbe/bot_agent_auth/application/exists/handler"
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
	user_register_action "github.com/bborbe/bot_agent_auth/user/register/action"
	user_register_handler "github.com/bborbe/bot_agent_auth/user/register/handler"
	user_remove_group_action "github.com/bborbe/bot_agent_auth/user/remove_group/action"
	user_remove_group_handler "github.com/bborbe/bot_agent_auth/user/remove_group/handler"
	user_unregister_action "github.com/bborbe/bot_agent_auth/user/unregister/action"
	user_unregister_handler "github.com/bborbe/bot_agent_auth/user/unregister/handler"
	user_whoami_action "github.com/bborbe/bot_agent_auth/user/whoami/action"
	user_whoami_handler "github.com/bborbe/bot_agent_auth/user/whoami/handler"
	flag "github.com/bborbe/flagenv"
	http_client_builder "github.com/bborbe/http/client_builder"
	"github.com/bborbe/http/header"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

const (
	PARAMETER_LOGLEVEL                  = "loglevel"
	PARAMETER_NSQ_LOOKUPD               = "nsq-lookupd-address"
	PARAMETER_NSQD                      = "nsqd-address"
	DEFAULT_BOT_NAME                    = "auth"
	PARAMETER_BOT_NAME                  = "bot-name"
	PARAMETER_ADMIN                     = "admin"
	PARAMETER_AUTH_ADDRESS              = "auth-address"
	PARAMETER_AUTH_APPLICATION_NAME     = "auth-application-name"
	PARAMETER_AUTH_APPLICATION_PASSWORD = "auth-application-password"
	PREFIX                              = "/auth"
	PARAMETER_RESTRICT_TO_TOKENS        = "restrict-to-tokens"
)

var (
	logger                     = log.DefaultLogger
	logLevelPtr                = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	nsqLookupdAddressPtr       = flag.String(PARAMETER_NSQ_LOOKUPD, "", "nsq lookupd address")
	nsqdAddressPtr             = flag.String(PARAMETER_NSQD, "", "nsqd address")
	botNamePtr                 = flag.String(PARAMETER_BOT_NAME, DEFAULT_BOT_NAME, "bot name")
	authAddressPtr             = flag.String(PARAMETER_AUTH_ADDRESS, "", "auth address")
	authApplicationNamePtr     = flag.String(PARAMETER_AUTH_APPLICATION_NAME, "", "auth application name")
	authApplicationPasswordPtr = flag.String(PARAMETER_AUTH_APPLICATION_PASSWORD, "", "auth application password")
	adminAuthTokenPtr          = flag.String(PARAMETER_ADMIN, "", "admin")
	restrictToTokensPtr        = flag.String(PARAMETER_RESTRICT_TO_TOKENS, "", "restrict to tokens")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)
	err := do(
		PREFIX,
		*nsqdAddressPtr,
		*nsqLookupdAddressPtr,
		*botNamePtr,
		*authAddressPtr,
		*authApplicationNamePtr,
		*authApplicationPasswordPtr,
		*adminAuthTokenPtr,
		*restrictToTokensPtr,
	)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(
	prefix string,
	nsqdAddress string,
	nsqLookupdAddress string,
	botname string,
	authAddress string,
	authApplicationName string,
	authApplicationPassword string,
	adminAuthToken string,
	restrictToTokens string,
) error {
	requestConsumer, err := createRequestConsumer(
		prefix,
		nsqdAddress,
		nsqLookupdAddress,
		botname,
		authAddress,
		authApplicationName,
		authApplicationPassword,
		adminAuthToken,
		restrictToTokens,
	)
	if err != nil {
		return err
	}
	return requestConsumer.Run()
}

func createRequestConsumer(
	prefix string,
	nsqdAddress string,
	nsqLookupdAddress string,
	botname string,
	authAddress string,
	authApplicationName string,
	authApplicationPassword string,
	adminAuthToken string,
	restrictToTokens string,
) (request_consumer.RequestConsumer, error) {
	if len(nsqLookupdAddress) == 0 {
		return nil, fmt.Errorf("parameter %s missing", PARAMETER_NSQ_LOOKUPD)
	}
	if len(nsqdAddress) == 0 {
		return nil, fmt.Errorf("parameter %s missing", PARAMETER_NSQD)
	}
	if len(botname) == 0 {
		return nil, fmt.Errorf("parameter %s missing", PARAMETER_BOT_NAME)
	}
	if len(authAddress) == 0 {
		return nil, fmt.Errorf("parameter %s missing", PARAMETER_AUTH_ADDRESS)
	}
	if len(authApplicationName) == 0 {
		return nil, fmt.Errorf("parameter %s missing", PARAMETER_AUTH_APPLICATION_NAME)
	}
	if len(authApplicationPassword) == 0 {
		return nil, fmt.Errorf("parameter %s missing", PARAMETER_AUTH_APPLICATION_PASSWORD)
	}

	httpRequestBuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	httpClient := http_client_builder.New().WithoutProxy().Build()

	restCaller := rest.New(authAddress, httpClient.Do, httpRequestBuilderProvider)

	token := header.CreateAuthorizationToken(authApplicationName, authApplicationPassword)

	applicationCreatorAction := application_creator_action.New(restCaller.Call, token)
	applicationCreatorHandler := application_creator_handler.New(prefix, adminAuthToken, applicationCreatorAction.Create)

	applicationDeletorAction := application_deletor_action.New(restCaller.Call, token)
	applicationDeletorHandler := application_deletor_handler.New(prefix, adminAuthToken, applicationDeletorAction.Delete)

	applicationExistsAction := application_exists_action.New(restCaller.Call, token)
	applicationExistsHandler := application_exists_handler.New(prefix, adminAuthToken, applicationExistsAction.Exists)

	userWhoamiAction := user_whoami_action.New(restCaller.Call, token)
	userWhoamiHandler := user_whoami_handler.New(prefix, userWhoamiAction.Whoami)

	userRegisterAction := user_register_action.New(restCaller.Call, token)
	userRegisterHandler := user_register_handler.New(prefix, userRegisterAction.Register)

	userUnregisterAction := user_unregister_action.New(restCaller.Call, token)
	userUnregisterHandler := user_unregister_handler.New(prefix, userUnregisterAction.Unregister)

	userCreateAction := user_create_action.New(restCaller.Call, token)
	userCreateHandler := user_create_handler.New(prefix, adminAuthToken, userCreateAction.CreateUser)

	userDeleteAction := user_delete_action.New(restCaller.Call, token)
	userDeleteHandler := user_delete_handler.New(prefix, adminAuthToken, userDeleteAction.DeleteUser)

	tokenAddAction := token_add_action.New(restCaller.Call, token)
	tokenAddHandler := token_add_handler.New(prefix, tokenAddAction.Add)

	tokenRemoveAction := token_remove_action.New(restCaller.Call, token)
	tokenRemoveHandler := token_remove_handler.New(prefix, tokenRemoveAction.Remove)

	userAddGroupAction := user_add_group_action.New(restCaller.Call, token)
	userAddGroupHandler := user_add_group_handler.New(prefix, adminAuthToken, userAddGroupAction.AddGroupToUser)

	userRemoveGroupAction := user_remove_group_action.New(restCaller.Call, token)
	userRemoveGroupHandler := user_remove_group_handler.New(prefix, adminAuthToken, userRemoveGroupAction.RemoveGroupToUser)

	producer, err := producer.New(nsqdAddress)
	if err != nil {
		return nil, err
	}

	sender := sender.New(producer)

	var messageHandler api.MessageHandler = match.New(
		prefix,
		applicationCreatorHandler,
		applicationDeletorHandler,
		applicationExistsHandler,
		userWhoamiHandler,
		userRegisterHandler,
		userUnregisterHandler,
		userCreateHandler,
		userDeleteHandler,
		tokenAddHandler,
		tokenRemoveHandler,
		userAddGroupHandler,
		userRemoveGroupHandler,
	)

	tokens := strings.Split(restrictToTokens, ",")
	if len(tokens) > 0 {
		messageHandler = restrict_to_tokens.New(
			messageHandler,
			tokens,
		)
	}

	return request_consumer.New(sender.Send, nsqdAddress, nsqLookupdAddress, botname, messageHandler), nil
}
