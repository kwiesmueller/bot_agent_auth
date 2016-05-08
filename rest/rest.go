package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bborbe/log"

	"time"

	"github.com/bborbe/http/bearer"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
)

var logger = log.DefaultLogger

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type rest struct {
	applicationName            string
	applicationPassword        string
	address                    string
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
	executeRequest             ExecuteRequest
}

type Rest interface {
	Call(path string, method string, request interface{}, response interface{}) error
}

func New(applicationName string, applicationPassword string, address string, executeRequest ExecuteRequest, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider) *rest {
	r := new(rest)
	r.applicationName = applicationName
	r.applicationPassword = applicationPassword
	r.address = address
	r.httpRequestBuilderProvider = httpRequestBuilderProvider
	r.executeRequest = executeRequest
	return r
}

func (r *rest) Call(path string, method string, request interface{}, response interface{}) error {
	logger.Debugf("call path %s on %s", path, r.applicationName)
	start := time.Now()
	defer logger.Debugf("create completed in %dms", time.Now().Sub(start)/time.Millisecond)
	target := fmt.Sprintf("http://%s%s", r.address, path)
	logger.Debugf("send message to %s", target)
	requestbuilder := r.httpRequestBuilderProvider.NewHttpRequestBuilder(target)
	requestbuilder.SetMethod(method)
	requestbuilder.AddContentType("application/json")
	requestbuilder.AddHeader("Authorization", bearer.CreateBearerHeader(r.applicationName, r.applicationPassword))
	content, err := json.Marshal(request)
	if err != nil {
		return err
	}
	logger.Debugf("send request to %s: %s", path, string(content))
	requestbuilder.SetBody(bytes.NewBuffer(content))
	req, err := requestbuilder.Build()
	if err != nil {
		return err
	}
	resp, err := r.executeRequest(req)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("request to %s failed with status: %d", path, resp.StatusCode)
	}
	if err = json.NewDecoder(resp.Body).Decode(response); err != nil {
		return err
	}
	return nil
}
