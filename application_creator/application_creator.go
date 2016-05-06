package application_creator

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

func (s *applicationCreator) Create(applicationName string) (*api.ApplicationPassword, error) {
	logger.Debugf("create application %s", applicationName)
	start := time.Now()
	defer logger.Debugf("create completed in %dms", time.Now().Sub(start)/time.Millisecond)
	target := fmt.Sprintf("http://%s/application", s.address)
	logger.Debugf("send message to %s", target)
	requestbuilder := s.httpRequestBuilderProvider.NewHttpRequestBuilder(target)
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	requestbuilder.AddHeader("Authorization", bearer.CreateBearerHeader(s.applicationName, s.applicationPassword))
	content, err := json.Marshal(&api.CreateApplicationRequest{
		ApplicationName: api.ApplicationName(applicationName),
	})
	if err != nil {
		return nil, err
	}
	logger.Debugf("send create application request to auth api: %s", string(content))
	requestbuilder.SetBody(bytes.NewBuffer(content))
	req, err := requestbuilder.Build()
	if err != nil {
		return nil, err
	}
	resp, err := s.executeRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("request to auth api failed with status: %d", resp.StatusCode)
	}
	var response api.CreateApplicationResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response.ApplicationPassword, nil
}
