package application_deletor

import (
	"fmt"
	"net/http"

	"github.com/bborbe/log"

	"time"

	"github.com/bborbe/http/bearer"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
)

var logger = log.DefaultLogger

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type applicationDeletor struct {
	applicationName            string
	applicationPassword        string
	address                    string
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
	executeRequest             ExecuteRequest
}

func New(applicationName string, applicationPassword string, address string, executeRequest ExecuteRequest, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider) *applicationDeletor {
	m := new(applicationDeletor)
	m.applicationName = applicationName
	m.applicationPassword = applicationPassword
	m.address = address
	m.httpRequestBuilderProvider = httpRequestBuilderProvider
	m.executeRequest = executeRequest
	return m
}

func (s *applicationDeletor) Delete(applicationName string) error {
	logger.Debugf("create application %s", applicationName)
	start := time.Now()
	defer logger.Debugf("create completed in %dms", time.Now().Sub(start)/time.Millisecond)
	target := fmt.Sprintf("http://%s/application/%s", s.address, applicationName)
	logger.Debugf("send message to %s", target)
	requestbuilder := s.httpRequestBuilderProvider.NewHttpRequestBuilder(target)
	requestbuilder.SetMethod("DELETE")
	requestbuilder.AddContentType("application/json")
	requestbuilder.AddHeader("Authorization", bearer.CreateBearerHeader("auth", s.applicationPassword))
	req, err := requestbuilder.Build()
	if err != nil {
		return err
	}
	resp, err := s.executeRequest(req)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}
	return nil
}
