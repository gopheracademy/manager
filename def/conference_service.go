package def

// ConferenceService is a service for managing Conferences
type ConferenceService interface {
	// Greet prepares a lovely greeting.
	List(ListConferenceRequest) ListConferenceResponse
	Get(GetConferenceRequest) GetConferenceResponse
}

// GetConferenceRequest is the request object for ConferenceService.Get.
type GetConferenceRequest struct {
}

// GetConferenceResponse is the response object containing a
// single Conference
type GetConferenceResponse struct {
	// Conference represents an event
	// like GopherCon 2020
	Conference Conference
}

// ListConferenceRequest is the request object for ConferenceService.List.
type ListConferenceRequest struct {
}

// ListConferenceResponse is the response object containing a
// list of Conferences
type ListConferenceResponse struct {
	// Greeting is a nice message welcoming somebody.
	Conferences []Conference
}
