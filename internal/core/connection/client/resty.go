package client

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-resty/resty/v2"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/generic"
)

var (
	contentTypeJson = "application/json"
)

// New new client
func New() *resty.Client {
	var debug bool
	if !config.CF.App.Environment.Production() {
		debug = true
	}

	client := resty.New()
	client.SetDebug(debug)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(3 * time.Minute)
	client.SetContentLength(true)

	return client
}

// Client client interface
type Client interface {
	BasicAuthentication(token string) string
	BearerAuthentication(token string) string
	NewHeaders(h map[string]string) map[string]string
	Get(url string, headers, queryParams map[string]string, i interface{}, breakerName string) (*resty.Response, error)
	Post(url string, headers, pathParams map[string]string, body interface{}, i interface{}, breakerName string) (*resty.Response, error)
	Put(url string, headers, pathParams map[string]string, body interface{}, i interface{}, breakerName string) (*resty.Response, error)
	Patch(url string, headers, pathParams map[string]string, body interface{}, i interface{}, breakerName string) (*resty.Response, error)
	Delete(url string, headers, pathParams map[string]string, body interface{}, i interface{}, breakerName string) (*resty.Response, error)
}

type client struct {
	session *resty.Client
}

// NewClient new client
func NewClient() Client {
	return &client{
		session: New(),
	}
}

// BasicAuthentication get basic token
func (c *client) BasicAuthentication(token string) string {
	return fmt.Sprintf("Basic %s", token)
}

// BearerAuthentication get bearer token
func (c *client) BearerAuthentication(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}

// NewHeaders new headers
func (c *client) NewHeaders(h map[string]string) map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = contentTypeJson
	headers["User-Agent"] = config.CF.App.ProjectName

	for k, v := range h {
		headers[k] = v
	}

	return headers
}

func (c *client) get(url string, headers, queryParams map[string]string, i interface{}) (*resty.Response, error) {
	req := c.session.
		R().
		SetHeaders(headers).
		SetQueryParams(queryParams)
	if i != nil {
		req = req.SetResult(i)
	}

	return req.Get(url)
}

func (c *client) post(url string, headers, pathParams map[string]string, body interface{}, i interface{}) (*resty.Response, error) {
	return c.setPost(headers, pathParams, body, i).Post(url)
}

func (c *client) put(url string, headers, pathParams map[string]string, body interface{}, i interface{}) (*resty.Response, error) {
	return c.setPost(headers, pathParams, body, i).Put(url)
}

func (c *client) patch(url string, headers, pathParams map[string]string, body interface{}, i interface{}) (*resty.Response, error) {
	return c.setPost(headers, pathParams, body, i).Patch(url)
}

func (c *client) delete(url string, headers, pathParams map[string]string, body interface{}, i interface{}) (*resty.Response, error) {
	return c.setPost(headers, pathParams, body, i).Delete(url)
}

func (c *client) setPost(headers, pathParams map[string]string, body interface{}, i interface{}) *resty.Request {
	req := c.session.
		R().
		SetHeaders(headers).
		SetPathParams(pathParams).
		SetBody(body)
	if i != nil {
		req = req.SetResult(i)
	}

	return req
}

// Breaker get request
func (c *client) Get(url string, headers, queryParams map[string]string, i interface{}, breakerName string) (*resty.Response, error) {
	if generic.IsEmpty(breakerName) {
		res, err := c.get(url, headers, queryParams, i)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	// breaker
	output := make(chan *resty.Response, 1)
	hystrix.Go(breakerName, func() error {
		res, err := c.get(url, headers, queryParams, i)
		if err != nil {
			return err
		}

		output <- res

		return nil
	}, nil)

	cb, _, _ := hystrix.GetCircuit(breakerName)
	if cb.IsOpen() {
		return nil, config.RR.Internal.TooManyRequests
	} else {
		return <-output, nil
	}
}

// Breaker post request
func (c *client) Post(url string, headers, pathParams map[string]string, body interface{}, i interface{}, breakerName string) (*resty.Response, error) {
	if generic.IsEmpty(breakerName) {
		res, err := c.post(url, headers, pathParams, body, i)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	// breaker
	output := make(chan *resty.Response, 1)
	hystrix.Go(breakerName, func() error {
		res, err := c.post(url, headers, pathParams, body, i)
		if err != nil {
			return err
		}

		output <- res

		return nil
	}, nil)

	cb, _, _ := hystrix.GetCircuit(breakerName)
	if cb.IsOpen() {
		return nil, config.RR.Internal.TooManyRequests
	} else {
		return <-output, nil
	}
}

// Breaker put request
func (c *client) Put(url string, headers, pathParams map[string]string, body interface{}, i interface{}, breakerName string) (*resty.Response, error) {
	if generic.IsEmpty(breakerName) {
		res, err := c.put(url, headers, pathParams, body, i)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	// breaker
	output := make(chan *resty.Response, 1)
	hystrix.Go(breakerName, func() error {
		res, err := c.put(url, headers, pathParams, body, i)
		if err != nil {
			return err
		}

		output <- res

		return nil
	}, nil)

	cb, _, _ := hystrix.GetCircuit(breakerName)
	if cb.IsOpen() {
		return nil, config.RR.Internal.TooManyRequests
	} else {
		return <-output, nil
	}
}

// Breaker patch request
func (c *client) Patch(url string, headers, pathParams map[string]string, body interface{}, i interface{}, breakerName string) (*resty.Response, error) {
	if generic.IsEmpty(breakerName) {
		res, err := c.patch(url, headers, pathParams, body, i)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	// breaker
	output := make(chan *resty.Response, 1)
	hystrix.Go(breakerName, func() error {
		res, err := c.patch(url, headers, pathParams, body, i)
		if err != nil {
			return err
		}

		output <- res

		return nil
	}, nil)

	cb, _, _ := hystrix.GetCircuit(breakerName)
	if cb.IsOpen() {
		return nil, config.RR.Internal.TooManyRequests
	} else {
		return <-output, nil
	}
}

// Breaker delete request
func (c *client) Delete(url string, headers, pathParams map[string]string, body interface{}, i interface{}, breakerName string) (*resty.Response, error) {
	if generic.IsEmpty(breakerName) {
		res, err := c.delete(url, headers, pathParams, body, i)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	// breaker
	output := make(chan *resty.Response, 1)
	hystrix.Go(breakerName, func() error {
		res, err := c.delete(url, headers, pathParams, body, i)
		if err != nil {
			return err
		}

		output <- res

		return nil
	}, nil)

	cb, _, _ := hystrix.GetCircuit(breakerName)
	if cb.IsOpen() {
		return nil, config.RR.Internal.TooManyRequests
	} else {
		return <-output, nil
	}
}
