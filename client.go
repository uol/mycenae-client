package mycenae

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/uol/funks"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//
// A client to mycenae timeseries project.
// author: rnojiri
//

// Client - a mycenae client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Configuration - client configurations
type Configuration struct {
	Host    string         `json:"host"`
	Port    int            `json:"port"`
	Secure  bool           `json:"secure"`
	Timeout funks.Duration `json:"timeout"`
}

// Validate - validates the configuration parameters
func (c *Configuration) Validate() error {

	if c == nil {
		return ErrNullConfiguration
	}

	if len(c.Host) == 0 {
		return ErrInvalidHost
	}

	if c.Port == 0 {
		return ErrInvalidPort
	}

	if c.Timeout.Duration == 0 {
		return ErrInvalidTimeout
	}

	return nil
}

type method string

type endpoint struct {
	uri              string
	method           string
	hasURIParameters bool
}

const (
	cHeaderKContentType string = "content-type"
	cHeaderVContentType string = "application/json"
	cMethodGET          string = "GET"
	cMethodPOST         string = "POST"
	cMethodHEAD         string = "HEAD"
	cMethodDELETE       string = "DELETE"
	cMethodPUT          string = "PUT"
)

var (
	// ErrBadRequest - raised when invalid parameters are passed
	ErrBadRequest error = errors.New("invalid parameters")

	// ErrNullConfiguration - raised when a null configuration is found
	ErrNullConfiguration error = errors.New("configuration is null")

	// ErrInvalidHost - raised when a host is invalid
	ErrInvalidHost error = errors.New("host is invalid")

	// ErrInvalidPort - raised when the port is invalid
	ErrInvalidPort error = errors.New("port is invalid")

	// ErrInvalidTimeout - raised when the timeout is invalid
	ErrInvalidTimeout error = errors.New("timeout is invalid")
)

// New - configures a new client
func New(configuration *Configuration) (*Client, error) {

	err := configuration.Validate()
	if err != nil {
		return nil, err
	}

	b := strings.Builder{}
	b.Grow(12 + len(configuration.Host))
	b.WriteString("http")

	if configuration.Secure {
		b.WriteString("s")
	}

	b.WriteString("://")
	b.WriteString(configuration.Host)

	if configuration.Port != 80 && configuration.Port != 443 {
		b.WriteString(":")
		b.WriteString(strconv.Itoa(configuration.Port))
	}

	return &Client{
		baseURL:    b.String(),
		httpClient: funks.CreateHTTPClient(configuration.Timeout.Duration, true),
	}, nil
}

// DoJSONRequest - creates a new GET request
func (c *Client) DoJSONRequest(e *endpoint, uriParameters []interface{}, body interface{}, resultJSON interface{}) (status int, err error) {

	var uri string
	if e.hasURIParameters {
		uri = fmt.Sprintf(e.uri, uriParameters...)
	} else {
		uri = e.uri
	}

	var iobody io.Reader
	if body != nil {
		var bodyBytes []byte
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return
		}

		iobody = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(e.method, fmt.Sprintf("%s%s", c.baseURL, uri), iobody)
	if err != nil {
		return
	}

	req.Header.Add(cHeaderKContentType, cHeaderVContentType)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return
	}

	defer func() {
		if res.Body != nil {
			res.Body.Close()
		}
	}()

	status = res.StatusCode

	if status != http.StatusNoContent && resultJSON != nil && res.Body != nil {
		var content []byte
		content, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}

		if len(content) > 0 {
			err = json.Unmarshal(content, &resultJSON)
			if err != nil {
				return
			}
		}
	}

	return
}

func newRequestError(status int) error {

	return fmt.Errorf("received an error or unmapped status: %d", status)
}

// Close - close idle connections
func (c *Client) Close() {

	c.httpClient.CloseIdleConnections()
}
