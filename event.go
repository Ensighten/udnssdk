package udnssdk

import (
	"fmt"
	"log"
	"time"
)

// EventsService manages Events
type EventsService struct {
	client *Client
}

// EventInfoDTO wraps an event's info response
type EventInfoDTO struct {
	ID         string    `json:"id"`
	PoolRecord string    `json:"poolRecord"`
	EventType  string    `json:"type"`
	Start      time.Time `json:"start"`
	Repeat     string    `json:"repeat"`
	End        time.Time `json:"end"`
	Notify     string    `json:"notify"`
}

// EventInfoListDTO wraps a list of event info and list metadata, from an index request
type EventInfoListDTO struct {
	Events     []EventInfoDTO `json:"events"`
	Queryinfo  QueryInfo      `json:"queryInfo"`
	Resultinfo ResultInfo     `json:"resultInfo"`
}

// EventPath generates a URI by zone, type, name & guid
func EventPath(zone, typ, name, guid string) string {
	e := &EventKey{Zone: zone, Type: typ, Name: name, GUID: guid}
	if guid == "" {
		return e.RRSetKey().EventsURI()
	}
	return e.URI()
}

// ListAllEvents requests all events, using pagination and error handling
func (s *SBTCService) ListAllEvents(query, name, typ, zone string) ([]EventInfoDTO, error) {
	r := RRSetKey{Zone: zone, Type: typ, Name: name}
	return s.Events().Select(r, query)
}

func eventQueryPath(zone, typ, name, query string, offset int) string {
	r := RRSetKey{Zone: zone, Type: typ, Name: name}
	return r.EventsQueryURI(query, offset)
}

// ListEvents requests list of events by query, name, type & zone, and offset, also returning list metadata, the actual response, or an error
func (s *SBTCService) ListEvents(query, name, typ, zone string, offset int) ([]EventInfoDTO, ResultInfo, *Response, error) {
	r := RRSetKey{Zone: zone, Type: typ, Name: name}
	return s.Events().SelectWithOffset(r, query, offset)
}

// GetEvent requests an event by name, type, zone & guid, also returning the actual response, or an error
func (s *SBTCService) GetEvent(name, typ, zone, guid string) (EventInfoDTO, *Response, error) {
	e := EventKey{Zone: zone, Type: typ, Name: name, GUID: guid}
	return s.Events().Find(e)
}

// CreateEvent requests creation of an event by name, type, zone, with provided event-info, returning actual response or an error
func (s *SBTCService) CreateEvent(name, typ, zone string, ev EventInfoDTO) (*Response, error) {
	r := RRSetKey{Zone: zone, Type: typ, Name: name}
	return s.Events().Create(r, ev)
}

// UpdateEvent requests update of an event by name, type, zone & guid, withprovided event-info, returning the actual response or an error
func (s *SBTCService) UpdateEvent(name, typ, zone, guid string, ev EventInfoDTO) (*Response, error) {
	e := EventKey{Zone: zone, Type: typ, Name: name, GUID: guid}
	return s.Events().Update(e, ev)
}

// DeleteEvent requests deletion of an event by name, type, zone & guid, returning the actual response or an error
func (s *SBTCService) DeleteEvent(name, typ, zone, guid string) (*Response, error) {
	e := EventKey{Zone: zone, Type: typ, Name: name, GUID: guid}
	return s.Events().Delete(e)
}

// <<<<

// Events builds an EventService from an SBTCService
func (s *SBTCService) Events() *EventsService {
	return &EventsService{client: s.client}
}

// EventKey collects the identifiers of an Event
type EventKey struct {
	Zone string
	Type string
	Name string
	GUID string
}

// RRSetKey generates the RRSetKey for the EventKey
func (p *EventKey) RRSetKey() *RRSetKey {
	return &RRSetKey{
		Zone: p.Zone,
		Type: p.Type,
		Name: p.Name,
	}
}

// URI generates the URI for a probe
func (p *EventKey) URI() string {
	return fmt.Sprintf("%s/%s", p.RRSetKey().EventsURI(), p.GUID)
}

// Select requests all events, using pagination and error handling
func (s *EventsService) Select(r RRSetKey, query string) ([]EventInfoDTO, error) {
	// TODO: Sane Configuration for timeouts / retries
	maxerrs := 5
	waittime := 5 * time.Second

	// init accumulators
	pis := []EventInfoDTO{}
	offset := 0
	errcnt := 0

	for {
		reqEvents, ri, res, err := s.SelectWithOffset(r, query, offset)
		if err != nil {
			if res.StatusCode >= 500 {
				errcnt = errcnt + 1
				if errcnt < maxerrs {
					time.Sleep(waittime)
					continue
				}
			}
			return pis, err
		}

		log.Printf("ResultInfo: %+v\n", ri)
		for _, pi := range reqEvents {
			pis = append(pis, pi)
		}
		if ri.ReturnedCount+ri.Offset >= ri.TotalCount {
			return pis, nil
		}
		offset = ri.ReturnedCount + ri.Offset
		continue
	}
}

// SelectWithOffset requests list of events by RRSetKey, query and offset, also returning list metadata, the actual response, or an error
func (s *EventsService) SelectWithOffset(r RRSetKey, query string, offset int) ([]EventInfoDTO, ResultInfo, *Response, error) {
	var tld EventInfoListDTO

	uri := r.EventsQueryURI(query, offset)
	res, err := s.client.get(uri, &tld)

	pis := []EventInfoDTO{}
	for _, pi := range tld.Events {
		pis = append(pis, pi)
	}
	return pis, tld.Resultinfo, res, err
}

// Create requests creation of an event by RRSetKey, with provided event-info, returning actual response or an error
func (s *EventsService) Create(r RRSetKey, ev EventInfoDTO) (*Response, error) {
	var ignored interface{}
	return s.client.post(r.EventsURI(), ev, &ignored)
}

// Find requests an event by name, type, zone & guid, also returning the actual response, or an error
func (s *EventsService) Find(e EventKey) (EventInfoDTO, *Response, error) {
	var t EventInfoDTO
	res, err := s.client.get(e.URI(), &t)
	return t, res, err
}

// Update requests update of an event by EventKey, withprovided event-info, returning the actual response or an error
func (s *EventsService) Update(e EventKey, ev EventInfoDTO) (*Response, error) {
	var ignored interface{}
	return s.client.put(e.URI(), ev, &ignored)
}

// Delete requests deletion of an event by EventKey, returning the actual response or an error
func (s *EventsService) Delete(e EventKey) (*Response, error) {
	return s.client.delete(e.URI(), nil)
}
