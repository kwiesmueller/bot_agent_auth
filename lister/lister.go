package lister

import (
	"fmt"
	"net/http"

	"github.com/bborbe/log"

	"encoding/json"

	"time"

	http_requestbuilder "github.com/bborbe/http/requestbuilder"
)

var logger = log.DefaultLogger

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type authCreator struct {
	address                    string
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
	executeRequest             ExecuteRequest
}

func New(address string, executeRequest ExecuteRequest, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider) *authCreator {
	m := new(authCreator)
	m.address = address
	m.httpRequestBuilderProvider = httpRequestBuilderProvider
	m.executeRequest = executeRequest
	return m
}

func (s *authCreator) List(authToken string) ([]string, error) {
	logger.Debugf("list")
	start := time.Now()
	defer logger.Debugf("list completed in %dms", time.Now().Sub(start)/time.Millisecond)
	target := fmt.Sprintf("http://%s/auth", s.address)
	logger.Debugf("send message to %s", target)
	requestbuilder := s.httpRequestBuilderProvider.NewHttpRequestBuilder(target)
	requestbuilder.SetMethod("GET")
	requestbuilder.AddContentType("application/json")
	requestbuilder.AddHeader("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req, err := requestbuilder.Build()
	if err != nil {
		return nil, err
	}
	resp, err := s.executeRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}
	var list []string
	if err = json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}
	return list, nil
}
