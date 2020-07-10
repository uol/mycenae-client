package client_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/uol/mycenae-shared/raw"

	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"
	gotesthttp "github.com/uol/gotest/http"
	"github.com/uol/mycenae-client"
)

func randomMetric() string {
	return strings.ReplaceAll(strings.ToLower(randomdata.FullName(randomdata.Female)), " ", ".")
}

func randomTags() map[string]string {

	tags := map[string]string{}

	for j := 0; j < randomdata.Number(3, 10); j++ {
		key := strings.ToLower(randomdata.FirstName(randomdata.RandomGender))
		value := strings.ToLower(randomdata.LastName())
		tags[key] = value
	}

	return tags
}

func generateRandomRawNumberResult() *raw.NumberQueryResults {

	numRes := randomdata.Number(1, 5)
	results := make([]raw.NumberPoints, numRes)
	for i := 0; i < numRes; i++ {
		results[i] = raw.NumberPoints{
			Metadata: raw.Metadata{
				Metric: randomMetric(),
				Tags:   randomTags(),
			},
			Values: []raw.NumberPoint{},
		}

		for j := 0; j < randomdata.Number(5, 10); j++ {
			results[i].Values = append(results[i].Values, raw.NumberPoint{
				Timestamp: time.Now().Unix(),
				Value:     float64(randomdata.Number(500000)),
			})
		}
	}

	return &raw.NumberQueryResults{
		Results: results,
		Total:   len(results),
	}
}

func generateRandomRawTextResult() *raw.TextQueryResults {

	numRes := randomdata.Number(1, 5)
	results := make([]raw.TextPoints, numRes)
	for i := 0; i < numRes; i++ {
		results[i] = raw.TextPoints{
			Metadata: raw.Metadata{
				Metric: randomMetric(),
				Tags:   randomTags(),
			},
			Texts: []raw.TextPoint{},
		}

		for j := 0; j < randomdata.Number(5, 10); j++ {
			results[i].Texts = append(results[i].Texts, raw.TextPoint{
				Timestamp: time.Now().Unix(),
				Text:      randomdata.City(),
			})
		}
	}

	return &raw.TextQueryResults{
		Results: results,
		Total:   len(results),
	}
}

func getRawQueryResponses(rd *mycenaeRandomData) (successResponses, alternateResponses, emptyResponses, errorResponses []gotesthttp.ResponseData) {

	successResponses = []gotesthttp.ResponseData{
		{
			RequestData: gotesthttp.RequestData{
				URI:    "/api/query/raw",
				Method: "POST",
				Body:   mustMarshalJSON(rd.rawNumberResults),
			},
			Status: http.StatusOK,
		},
	}

	alternateResponses = []gotesthttp.ResponseData{
		{
			RequestData: gotesthttp.RequestData{
				URI:    "/api/query/raw",
				Method: "POST",
				Body:   mustMarshalJSON(rd.rawTextResults),
			},
			Status: http.StatusOK,
		},
	}

	emptyResponses = []gotesthttp.ResponseData{
		{
			RequestData: gotesthttp.RequestData{
				URI:    "/api/query/raw",
				Method: "POST",
			},
			Status: http.StatusNoContent,
		},
	}

	errorResponses = []gotesthttp.ResponseData{
		{
			RequestData: gotesthttp.RequestData{
				URI:    "/api/query/raw",
				Method: "POST",
			},
			Status: http.StatusInternalServerError,
		},
	}

	return
}

func buildRandomRawQuery(ctype string) *raw.Query {

	return &raw.Query{
		Metadata: raw.Metadata{
			Metric: randomMetric(),
			Tags:   randomTags(),
		},
		Since: fmt.Sprintf("%dm", randomdata.Number(1, 59)),
		Until: fmt.Sprintf("%dm", randomdata.Number(1, 10)),
		Type:  ctype,
	}
}

func testNumberRawQuery(t *testing.T, client *mycenae.Client, server *gotesthttp.Server, rd *mycenaeRandomData) {

	server.SetMode(cSuccessMode)

	query := buildRandomRawQuery("meta")
	jsonQuery := mustMarshalJSON(query)

	results, err := client.GetRawPoints(query)
	if !assert.NoError(t, err, "expected no error calling api") {
		return
	}

	if !assert.NotNil(t, results, "expected a not null result") {
		return
	}

	if !assert.Equal(t, results, rd.rawNumberResults, "expected same object") {
		return
	}

	checkRequestDetails(t, server, jsonQuery, "POST", "/api/query/raw")

	server.SetMode(cEmptyMode)

	results, err = client.GetRawPoints(query)
	if !assert.NoError(t, err, "expected no error calling api") {
		return
	}

	if !assert.Nil(t, results, "expected a null result") {
		return
	}

	checkRequestDetails(t, server, jsonQuery, "POST", "/api/query/raw")

	server.SetMode(cErrorMode)

	results, err = client.GetRawPoints(query)
	if !assert.Error(t, err, "expected error calling api") {
		return
	}

	if !assert.Nil(t, results, "expected a null result") {
		return
	}

	checkRequestDetails(t, server, jsonQuery, "POST", "/api/query/raw")
}

func testTextRawQuery(t *testing.T, client *mycenae.Client, server *gotesthttp.Server, rd *mycenaeRandomData) {

	server.SetMode(cAlternateMode)

	query := buildRandomRawQuery("metatext")
	jsonQuery := mustMarshalJSON(query)

	results, err := client.GetRawTextPoints(query)
	if !assert.NoError(t, err, "expected no error calling api") {
		return
	}

	if !assert.NotNil(t, results, "expected a not null result") {
		return
	}

	if !assert.Equal(t, results, rd.rawTextResults, "expected same object") {
		return
	}

	checkRequestDetails(t, server, jsonQuery, "POST", "/api/query/raw")

	server.SetMode(cEmptyMode)

	results, err = client.GetRawTextPoints(query)
	if !assert.NoError(t, err, "expected no error calling api") {
		return
	}

	if !assert.Nil(t, results, "expected a null result") {
		return
	}

	checkRequestDetails(t, server, jsonQuery, "POST", "/api/query/raw")

	server.SetMode(cErrorMode)

	results, err = client.GetRawTextPoints(query)
	if !assert.Error(t, err, "expected error calling api") {
		return
	}

	if !assert.Nil(t, results, "expected a null result") {
		return
	}

	checkRequestDetails(t, server, jsonQuery, "POST", "/api/query/raw")
}
