package client_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/uol/mycenae-shared/raw"

	"github.com/stretchr/testify/assert"

	"github.com/uol/funks"
	gotesthttp "github.com/uol/gotest/http"
	"github.com/uol/mycenae-client"
)

//
// The main test suite to reutilize the timeseries mock.
// author: rnojiri
//

const (
	host           string = "localhost"
	port           int    = 18080
	cSuccessMode   string = "success"
	cAlternateMode string = "alternate"
	cEmptyMode     string = "empty"
	cErrorMode     string = "error"
)

type mycenaeRandomData struct {
	keysets          []string
	rawNumberResults *raw.NumberQueryResults
	rawTextResults   *raw.TextQueryResults
}

func mustMarshalJSON(data interface{}) string {

	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func generateRandomData() *mycenaeRandomData {

	mr := &mycenaeRandomData{
		keysets:          generateRandomKeysets(),
		rawNumberResults: generateRandomRawNumberResult(),
		rawTextResults:   generateRandomRawTextResult(),
	}

	return mr
}

func setupMycenaeMock(rd *mycenaeRandomData) *gotesthttp.Server {

	successResponses := []gotesthttp.ResponseData{}
	alternateResponses := []gotesthttp.ResponseData{}
	emptyResponses := []gotesthttp.ResponseData{}
	errorResponses := []gotesthttp.ResponseData{}

	suc, alt, emp, err := getKeysetResponses(rd)
	successResponses = append(successResponses, suc...)
	alternateResponses = append(alternateResponses, alt...)
	emptyResponses = append(emptyResponses, emp...)
	errorResponses = append(errorResponses, err...)

	suc, alt, emp, err = getRawQueryResponses(rd)
	successResponses = append(successResponses, suc...)
	alternateResponses = append(alternateResponses, alt...)
	emptyResponses = append(emptyResponses, emp...)
	errorResponses = append(errorResponses, err...)

	conf := gotesthttp.Configuration{
		Host:        host,
		Port:        port,
		ChannelSize: 10,
		Responses: map[string][]gotesthttp.ResponseData{
			cSuccessMode:   successResponses,
			cAlternateMode: alternateResponses,
			cEmptyMode:     emptyResponses,
			cErrorMode:     errorResponses,
		},
	}

	return gotesthttp.NewServer(&conf)
}

func createClient() *mycenae.Client {

	conf := mycenae.Configuration{
		Host:    host,
		Port:    port,
		Secure:  false,
		Timeout: *funks.ForceNewStringDuration("3s"),
	}

	client, err := mycenae.New(&conf)
	if err != nil {
		panic(err)
	}

	return client
}

func TestSuite(t *testing.T) {

	mr := generateRandomData()

	server := setupMycenaeMock(mr)
	defer server.Close()

	client := createClient()

	tests := []struct {
		title string
		run   func(t *testing.T, client *mycenae.Client, server *gotesthttp.Server, rd *mycenaeRandomData)
	}{
		{"get keysets", testGetKeysets},
		{"check keyset", testKeysetExists},
		{"raw query number", testNumberRawQuery},
		{"raw text query", testTextRawQuery},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.run(t, client, server, mr)
		})
	}
}

func checkRequestDetails(t *testing.T, server *gotesthttp.Server, body, method, uri string) bool {

	reqData := gotesthttp.WaitForServerRequest(server, 1*time.Second, 3*time.Second)
	if !assert.Equal(t, reqData.Body, body, "expected an empty request body") {
		return false
	}

	if !assert.Equal(t, reqData.Headers.Get("content-type"), "application/json", "expected application/json as content-type") {
		return false
	}

	if !assert.Equal(t, reqData.Method, method, "expected GET as method") {
		return false
	}

	return assert.Equal(t, reqData.URI, uri, "expected same uri")
}
