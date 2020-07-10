package mycenae

import (
	"net/http"

	"github.com/uol/funks"
)

var (
	endpointRawGET endpoint = endpoint{method: cMethodPOST, uri: "/api/query/raw", hasURIParameters: false}
)

const (
	cNumberResult string = "meta"
	cTextResult   string = "metatext"
)

// RawDataQuery - the raw data query JSON
type RawDataQuery struct {
	Metric string            `json:"metric"`
	Tags   map[string]string `json:"tags"`
	Since  funks.Duration    `json:"since"`
	Until  funks.Duration    `json:"until"`
}

// RawDataQueryJSON - the raw data query JSON
type RawDataQueryJSON struct {
	RawDataQuery
	Type string `json:"type"`
}

// RawDataMetadata - the raw data (metadata only)
type RawDataMetadata struct {
	Metric string            `json:"metric"`
	Tags   map[string]string `json:"tags"`
}

// RawDataNumberPoint - represents a raw number point result
type RawDataNumberPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// RawDataTextPoint - represents a raw text point result
type RawDataTextPoint struct {
	Timestamp int64  `json:"timestamp"`
	Text      string `json:"text"`
}

// RawDataQueryNumberPoints - the metadata and value results
type RawDataQueryNumberPoints struct {
	Metadata RawDataMetadata      `json:"metadata"`
	Values   []RawDataNumberPoint `json:"points"`
}

// RawDataQueryTextPoints - the metadata and text results
type RawDataQueryTextPoints struct {
	Metadata RawDataMetadata    `json:"metadata"`
	Texts    []RawDataTextPoint `json:"points"`
}

// RawDataQueryNumberResults - the final raw query number results
type RawDataQueryNumberResults struct {
	Results []RawDataQueryNumberPoints `json:"results"`
	Total   int                        `json:"total"`
}

// RawDataQueryTextResults - the final raw query text results
type RawDataQueryTextResults struct {
	Results []RawDataQueryTextPoints `json:"results"`
	Total   int                      `json:"total"`
}

// GetRawPoints - return the list of raw points
func (c *Client) GetRawPoints(query *RawDataQuery) (*RawDataQueryNumberResults, error) {

	results := &RawDataQueryNumberResults{}

	err := c.getRawPoints(query, cNumberResult, results)
	if err != nil {
		return nil, err
	}

	if results.Total == 0 {
		return nil, nil
	}

	return results, nil
}

// GetRawTextPoints - return the list of raw text points
func (c *Client) GetRawTextPoints(query *RawDataQuery) (*RawDataQueryTextResults, error) {

	results := &RawDataQueryTextResults{}

	err := c.getRawPoints(query, cTextResult, results)
	if err != nil {
		return nil, err
	}

	if results.Total == 0 {
		return nil, nil
	}

	return results, nil
}

// getRawPoints - return the list of raw points
func (c *Client) getRawPoints(query *RawDataQuery, ctype string, results interface{}) error {

	jsonQuery := RawDataQueryJSON{
		RawDataQuery: *query,
		Type:         ctype,
	}

	status, err := c.DoJSONRequest(&endpointRawGET, nil, jsonQuery, results)
	if err != nil {
		return err
	}

	if status == http.StatusOK {

		return nil
	}

	if status == http.StatusNoContent {

		return nil
	}

	if status == http.StatusBadRequest {

		return ErrBadRequest
	}

	return newRequestError(status)
}
