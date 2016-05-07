package main

import (
	"fmt"
	"os"

	"github.com/bborbe/bot_agent/request_consumer"
	application_creator_action "github.com/bborbe/bot_agent_auth/application/creator/action"
	application_deletor_action "github.com/bborbe/bot_agent_auth/application/deletor/action"
	application_exists_action "github.com/bborbe/bot_agent_auth/application/exists/action"
	application_creator_handler "github.com/bborbe/bot_agent_auth/application/creator/handler"
	application_deletor_handler "github.com/bborbe/bot_agent_auth/application/deletor/handler"
	application_exists_handler "github.com/bborbe/bot_agent_auth/application/exists/handler"
	"github.com/bborbe/bot_agent_auth/message_handler"
	flag "github.com/bborbe/flagenv"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

const (
	PARAMETER_LOGLEVEL = "loglevel"
	PARAMETER_NSQ_LOOKUPD = "nsq-lookupd-address"
	PARAMETER_NSQD = "nsqd-address"
	DEFAULT_BOT_NAME = "auth"
	PARAMETER_BOT_NAME = "bot-name"
	PARAMETER_AUTH_ADDRESS = "auth-address"
	PARAMETER_AUTH_APPLICATION_NAME = "auth-application-name"
	PARAMETER_AUTH_APPLICATION_PASSWORD = "auth-application-password"
	PREFIX = "/auth"
)

var (
	logger = log.DefaultLogger
	logLevelPtr = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	nsqLookupdAddressPtr = flag.String(PARAMETER_NSQ_LOOKUPD, "", "nsq lookupd address")
	nsqdAddressPtr = flag.String(PARAMETER_NSQD, "", "nsqd address")
	botNamePtr = flag.String(PARAMETER_BOT_NAME, DEFAULT_BOT_NAME, "bot name")
	authAddressPtr = flag.String(PARAMETER_AUTH_ADDRESS, "", "auth address")
	authApplicationNamePtr = flag.String(PARAMETER_AUTH_APPLICATION_NAME, "", "auth application name")
	authApplicationPasswordPtr = flag.String(PARAMETER_AUTH_APPLICATION_PASSWORD, "", "auth application password")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)
	err := do(PREFIX, *nsqdAddressPtr, *nsqLookupdAddressPtr, *botNamePtr, *authAddressPtr, *authApplicationNamePtr, *authApplicationPasswordPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(prefix string, nsqdAddress string, nsqLookupdAddress string, botname string, authAddress string, authApplicationName string, authApplicationPassword string) error {
	requestConsumer, err := createRequestConsumer(prefix, nsqdAddress, nsqLookupdAddress, botname, authAddress, authApplicationName, authApplicationPassword)
	if err != nil {
		return err
	}
	return requestConsumer.Run()
}

func createRequestConsumer(prefix string, nsqdAddress string, nsqLookupdAddress string, botname string, authAddress string, authApplicationName string, authApplicationPassword string) (request_consumer.RequestConsumer, error) {
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
	applicationCreatorAction := application_creator_action.New(authApplicationName, authApplicationPassword, authAddress, httpClient.Do, httpRequestBuilderProvider)
	applicationDeletorAction := application_deletor_action.New(authApplicationName, authApplicationPassword, authAddress, httpClient.Do, httpRequestBuilderProvider)
	applicationExistsAction := application_exists_action.New(authApplicationName, authApplicationPassword, authAddress, httpClient.Do, httpRequestBuilderProvider)

	applicationCreatorHandler := application_creator_handler.New(prefix, applicationCreatorAction.Create)
	applicationDeletorHandler := application_deletor_handler.New(prefix, applicationDeletorAction.Delete)
	applicationExistsHandler := application_exists_handler.New(prefix, applicationExistsAction.Exists)

	messageHandler := message_handler.New(prefix, applicationCreatorHandler, applicationDeletorHandler, applicationExistsHandler)
	return request_consumer.New(nsqdAddress, nsqLookupdAddress, botname, messageHandler), nil
}
