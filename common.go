package udnssdk

import (
	"fmt"
)

// GetResultByURI just requests a URI
func (c *Client) GetResultByURI(uri string) (*Response, error) {
	req, err := c.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return &Response{Response: res}, err
	}
	return &Response{Response: res}, err
}

// RRSetKey collects the identifiers of a Zone
type RRSetKey struct {
	Zone string
	Type string
	Name string
}

// URI generates the URI for an RRSet
func (r *RRSetKey) URI() string {
	return fmt.Sprintf("zones/%s/rrsets/%s/%s", r.Zone, r.Type, r.Name)
}

// AlertsURI generates the URI for an RRSet
func (r *RRSetKey) AlertsURI() string {
	return fmt.Sprintf("%s/alerts", r.URI())
}

// AlertsQueryURI generates the alerts query URI for an RRSet with query
func (r *RRSetKey) AlertsQueryURI(offset int) string {
	uri := r.AlertsURI()
	if offset != 0 {
		uri = fmt.Sprintf("%s?offset=%d", uri, offset)
	}
	return uri
}

// ProbesURI generates the probes URI for an RRSet
func (r *RRSetKey) ProbesURI() string {
	return fmt.Sprintf("%s/probes", r.URI())
}

// ProbesQueryURI generates the probes query URI for an RRSet with query
func (r *RRSetKey) ProbesQueryURI(query string) string {
	uri := r.ProbesURI()
	if query != "" {
		uri = fmt.Sprintf("%s?sort=NAME&query=%s", uri, query)
	}
	return uri
}
