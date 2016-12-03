package main

import (
	"runtime"

	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/bot_agent_auth/factory"
	"github.com/bborbe/bot_agent_auth/model"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/nsq_utils"
	"github.com/bborbe/nsq_utils/producer"
	"github.com/golang/glog"
)

const (
	parameterNsqLookupd              = "nsq-lookupd-address"
	parameterNsqd                    = "nsqd-address"
	parameterBotName                 = "bot-name"
	parameterAuthUrl                 = "auth-url"
	parameterAuthApplicationName     = "auth-application-name"
	parameterAuthApplicationPassword = "auth-application-password"
	parameterRestrictToTokens        = "restrict-to-tokens"
	parameterPrefix                  = "prefix"
	parameterAdminTokens             = "admin-to-tokens"
	parameterAdminGroups             = "admin-to-groups"
)

var (
	nsqLookupdAddressPtr       = flag.String(parameterNsqLookupd, "", "nsq lookupd address")
	nsqdAddressPtr             = flag.String(parameterNsqd, "", "nsqd address")
	botNamePtr                 = flag.String(parameterBotName, "auth", "bot name")
	authUrlPtr                 = flag.String(parameterAuthUrl, "", "auth url")
	authApplicationNamePtr     = flag.String(parameterAuthApplicationName, "", "auth application name")
	authApplicationPasswordPtr = flag.String(parameterAuthApplicationPassword, "", "auth application password")
	prefixPtr                  = flag.String(parameterPrefix, "/auth", "prefix commands start with")
	restrictToTokensPtr        = flag.String(parameterRestrictToTokens, "", "restrict to tokens")
	adminAuthTokensPtr         = flag.String(parameterAdminTokens, "", "admin tokens")
	adminAuthGroupsPtr         = flag.String(parameterAdminGroups, "", "admin groups")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := do(); err != nil {
		glog.Exit(err)
	}
}

func do() error {
	config := createConfig()
	if err := config.Validate(); err != nil {
		return err
	}
	producer, err := producer.New(config.NsqdAddress)
	if err != nil {
		return err
	}
	factory := factory.New(config, producer)
	return factory.RequestConsumer().Run()
}

func createConfig() model.Config {
	return model.Config{
		Prefix:                  model.Prefix(*prefixPtr),
		NsqdAddress:             nsq_utils.NsqdAddress(*nsqdAddressPtr),
		NsqLookupdAddress:       nsq_utils.NsqLookupdAddress(*nsqLookupdAddressPtr),
		Botname:                 nsq_utils.NsqChannel(*botNamePtr),
		AuthUrl:                 auth_model.Url(*authUrlPtr),
		AuthApplicationName:     auth_model.ApplicationName(*authApplicationNamePtr),
		AuthApplicationPassword: auth_model.ApplicationPassword(*authApplicationPasswordPtr),
		RestrictToTokens:        auth_model.ParseTokens(*restrictToTokensPtr),
		AdminAuthTokens:         auth_model.ParseTokens(*adminAuthTokensPtr),
		AdminGroups:              auth_model.ParseGroupNames(*adminAuthGroupsPtr),
	}
}
