package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bborbe/log"

	"time"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/http/bearer"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
)

var logger = log.DefaultLogger

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type applicationCreator struct {
	applicationName            string
	applicationPassword        string
	address                    string
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
	executeRequest             ExecuteRequest
}

func New(applicationName string, applicationPassword string, address string, executeRequest ExecuteRequest, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider) *applicationCreator {
	m := new(applicationCreator)
	m.applicationName = applicationName
	m.applicationPassword = applicationPassword
	m.address = address
	m.httpRequestBuilderProvider = httpRequestBuilderProvider
	m.executeRequest = executeRequest
	return m
}

func (a *applicationCreator) Remove(authToken string, token string) error {
	logger.Debugf("remove token %s to user with token %s", token, authToken)
	start := time.Now()
	defer logger.Debugf("remove token completed in %dms", time.Now().Sub(start)/time.Millisecond)

	if authToken == token {
		return fmt.Errorf("token equals authToken")
	}

	target := fmt.Sprintf("http://%s/token", a.address)
	logger.Debugf("send message to %s", target)
	requestbuilder := a.httpRequestBuilderProvider.NewHttpRequestBuilder(target)
	requestbuilder.SetMethod("DELETE")
	requestbuilder.AddContentType("application/json")
	requestbuilder.AddHeader("Authorization", bearer.CreateBearerHeader(a.applicationName, a.applicationPassword))
	content, err := json.Marshal(&api.RemoveTokenRequest{
		AuthToken: api.AuthToken(authToken),
		Token:     api.AuthToken(token),
	})
	if err != nil {
		return err
	}
	logger.Debugf("send remove token request to auth api: %s", string(content))
	requestbuilder.SetBody(bytes.NewBuffer(content))
	req, err := requestbuilder.Build()
	if err != nil {
		return err
	}
	resp, err := a.executeRequest(req)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("remove token request to auth api failed with status: %d", resp.StatusCode)
	}
	return nil
}
