package mycenae

import (
	"net/http"
)

var (
	endpointKeysetsGET  endpoint = endpoint{method: cMethodGET, uri: "/keysets", hasURIParameters: false}
	endpointKeysetsHEAD endpoint = endpoint{method: cMethodHEAD, uri: "/keyset/%s", hasURIParameters: true}
)

// GetKeysets - return the list of all keysets
func (c *Client) GetKeysets() ([]string, error) {

	results := []string{}

	status, err := c.DoJSONRequest(&endpointKeysetsGET, nil, nil, &results)
	if err != nil {
		return nil, err
	}

	if status == http.StatusOK || status == http.StatusNoContent {

		return results, nil
	}

	return nil, newRequestError(status)
}

// KeysetExists - check if a keyset exists
func (c *Client) KeysetExists(keyset string) (bool, error) {

	status, err := c.DoJSONRequest(&endpointKeysetsHEAD, []interface{}{keyset}, nil, nil)
	if err != nil {
		return false, err
	}

	if status == http.StatusOK {

		return true, nil
	}

	if status == http.StatusNotFound {

		return false, nil
	}

	return false, newRequestError(status)
}
