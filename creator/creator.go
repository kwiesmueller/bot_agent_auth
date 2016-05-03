package creator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bborbe/log"

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

type request struct {
	Name string `json:"name"`
}

func (s *authCreator) Create(authToken string, name string) error {
	logger.Debugf("create %s", name)
	start := time.Now()
	defer logger.Debugf("create completed in %dms", time.Now().Sub(start)/time.Millisecond)
	target := fmt.Sprintf("http://%s/auth", s.address)
	logger.Debugf("send message to %s", target)
	requestbuilder := s.httpRequestBuilderProvider.NewHttpRequestBuilder(target)
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	requestbuilder.AddHeader("Authorization", fmt.Sprintf("Bearer %s", authToken))
	content, err := json.Marshal(&request{
		Name: name,
	})
	if err != nil {
		return err
	}
	logger.Debugf("request message: %s", string(content))
	requestbuilder.SetBody(bytes.NewBuffer(content))
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
