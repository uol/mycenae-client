package client_test

import (
	"net/http"
	"strings"
	"testing"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"
	gotesthttp "github.com/uol/gotest/http"
	"github.com/uol/mycenae-client"
)

func generateRandomKeysets() []string {

	keysets := make([]string, randomdata.Number(5, 10))

	for i := 0; i < len(keysets); i++ {
		keysets[i] = strings.ToLower(randomdata.SillyName())
	}

	return keysets
}

func getKeysetResponses(rd *mycenaeRandomData) (successResponses, alternateResponses, emptyResponses, errorResponses []gotesthttp.ResponseData) {

	successResponses = []gotesthttp.ResponseData{
		{
			RequestData: gotesthttp.RequestData{
				URI:    "/keysets",
				Method: "GET",
				Body:   mustMarshalJSON(rd.keysets),
			},
			Status: http.StatusOK,
		},
		{
			RequestData: gotesthttp.RequestData{
				URI:    "/keyset/" + rd.keysets[len(rd.keysets)-1],
				Method: "HEAD",
			},
			Status: http.StatusOK,
		},
	}

	emptyResponses = []gotesthttp.ResponseData{
		{
			RequestData: gotesthttp.RequestData{
				URI:    "/keysets",
				Method: "GET",
			},
			Status: http.StatusNoContent,
		},
		{
			RequestData: gotesthttp.RequestData{
				URI:    "/keyset/" + rd.keysets[0],
				Method: "HEAD",
			},
			Status: http.StatusNotFound,
		},
	}

	errorResponses = []gotesthttp.ResponseData{
		{
			RequestData: gotesthttp.RequestData{
				URI:    "/keysets",
				Method: "GET",
			},
			Status: http.StatusInternalServerError,
		},
		{
			RequestData: gotesthttp.RequestData{
				URI:    "/keyset/" + rd.keysets[0],
				Method: "HEAD",
			},
			Status: http.StatusInternalServerError,
		},
	}

	return
}

func testGetKeysets(t *testing.T, client *mycenae.Client, server *gotesthttp.Server, rd *mycenaeRandomData) {

	server.SetMode(cSuccessMode)

	keysets, err := client.GetKeysets()
	if !assert.NoError(t, err, "expected no error calling api") {
		return
	}

	if !assert.Len(t, keysets, len(rd.keysets), "expected same length of keysets") {
		return
	}

	for i := 0; i < len(rd.keysets); i++ {
		if !assert.Equal(t, rd.keysets[i], keysets[i], "expected same keyset") {
			return
		}
	}

	checkRequestDetails(t, server, "", "GET", "/keysets")

	server.SetMode(cEmptyMode)

	keysets, err = client.GetKeysets()
	if !assert.NoError(t, err, "expected no error calling api") {
		return
	}

	if !assert.Len(t, keysets, 0, "expected no keysets") {
		return
	}

	checkRequestDetails(t, server, "", "GET", "/keysets")

	server.SetMode(cErrorMode)

	keysets, err = client.GetKeysets()
	if !assert.Error(t, err, "expected error calling api") {
		return
	}

	if !assert.Len(t, keysets, 0, "expected no keysets") {
		return
	}

	checkRequestDetails(t, server, "", "GET", "/keysets")
}

func testKeysetExists(t *testing.T, client *mycenae.Client, server *gotesthttp.Server, rd *mycenaeRandomData) {

	server.SetMode(cSuccessMode)

	exists, err := client.KeysetExists(rd.keysets[len(rd.keysets)-1])
	if !assert.NoError(t, err, "expected no error calling api") {
		return
	}

	if !assert.True(t, exists, "expected true as response") {
		return
	}

	if !checkRequestDetails(t, server, "", "HEAD", "/keyset/"+rd.keysets[len(rd.keysets)-1]) {
		return
	}

	server.SetMode(cEmptyMode)

	exists, err = client.KeysetExists(rd.keysets[0])
	if !assert.NoError(t, err, "expected no error calling api") {
		return
	}

	if !assert.False(t, exists, "expected false as response") {
		return
	}

	checkRequestDetails(t, server, "", "HEAD", "/keyset/"+rd.keysets[0])

	server.SetMode(cErrorMode)

	exists, err = client.KeysetExists(rd.keysets[0])
	if !assert.Error(t, err, "expected error calling api") {
		return
	}

	if !assert.False(t, exists, "expected false as response") {
		return
	}

	checkRequestDetails(t, server, "", "HEAD", "/keyset/"+rd.keysets[0])
}
