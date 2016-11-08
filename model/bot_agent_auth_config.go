package model

import (
	"fmt"

	auth_model "github.com/bborbe/auth/model"
	"github.com/bborbe/nsq_utils"
)

type Config struct {
	Prefix                  Prefix
	NsqdAddress             nsq_utils.NsqdAddress
	NsqLookupdAddress       nsq_utils.NsqLookupdAddress
	Botname                 nsq_utils.NsqChannel
	AuthUrl                 auth_model.Url
	AuthApplicationName     auth_model.ApplicationName
	AuthApplicationPassword auth_model.ApplicationPassword
	AdminAuthToken          auth_model.AuthToken
	RestrictToTokens        []auth_model.AuthToken
}

func (c *Config) Validate() error {
	if len(c.Prefix) == 0 {
		return fmt.Errorf("parameter Prefix missing")
	}
	if len(c.NsqdAddress) == 0 {
		return fmt.Errorf("parameter NsqdAddress missing")
	}
	if len(c.NsqLookupdAddress) == 0 {
		return fmt.Errorf("parameter NsqLookupdAddress missing")
	}
	if len(c.Botname) == 0 {
		return fmt.Errorf("parameter Botname missing")
	}
	if len(c.AuthUrl) == 0 {
		return fmt.Errorf("parameter AuthUrl missing")
	}
	if len(c.AuthApplicationName) == 0 {
		return fmt.Errorf("parameter AuthApplicationName missing")
	}
	if len(c.AuthApplicationPassword) == 0 {
		return fmt.Errorf("parameter AuthApplicationPassword missing")
	}
	return nil
}
