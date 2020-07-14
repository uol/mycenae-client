package mycenae

import (
	"net/http"

	"github.com/uol/mycenae-shared/raw"
)

var (
	endpointRawGET endpoint = endpoint{method: cMethodPOST, uri: "/api/query/raw", hasURIParameters: false}
)

const (
	cNumberResult string = "number"
	cTextResult   string = "text"
)

// GetRawPoints - return the list of raw points
func (c *Client) GetRawPoints(query *raw.Query) (*raw.NumberQueryResults, error) {

	results := &raw.NumberQueryResults{}

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
func (c *Client) GetRawTextPoints(query *raw.Query) (*raw.TextQueryResults, error) {

	results := &raw.TextQueryResults{}

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
func (c *Client) getRawPoints(query *raw.Query, ctype string, results interface{}) error {

	query.Type = ctype

	status, err := c.DoJSONRequest(&endpointRawGET, nil, query, results)
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
