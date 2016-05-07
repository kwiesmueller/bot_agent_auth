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

func (a *applicationCreator) Whoami(authToken string) (*api.UserName, error) {
	logger.Debugf("who is %s", authToken)
	start := time.Now()
	defer logger.Debugf("whoami completed in %dms", time.Now().Sub(start)/time.Millisecond)
	target := fmt.Sprintf("http://%s/login", a.address)
	logger.Debugf("send message to %s", target)
	requestbuilder := a.httpRequestBuilderProvider.NewHttpRequestBuilder(target)
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	requestbuilder.AddHeader("Authorization", bearer.CreateBearerHeader(a.applicationName, a.applicationPassword))
	content, err := json.Marshal(&api.LoginRequest{
		AuthToken: api.AuthToken(authToken),
	})
	if err != nil {
		return nil, err
	}
	logger.Debugf("send login request to auth api: %s", string(content))
	requestbuilder.SetBody(bytes.NewBuffer(content))
	req, err := requestbuilder.Build()
	if err != nil {
		return nil, err
	}
	resp, err := a.executeRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("request to auth api failed with status: %d", resp.StatusCode)
	}
	var response api.LoginResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response.User, nil
}
