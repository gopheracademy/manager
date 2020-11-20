package def

// ConferenceService is a service for managing Conferences
type ConferenceService interface {
	// Greet prepares a lovely greeting.
	List(ListConferenceRequest) ListConferenceResponse
	Get(GetConferenceRequest) GetConferenceResponse
	GetBySlug(GetConferenceBySlugRequest) GetConferenceResponse
	Create(CreateConferenceRequest) CreateConferenceResponse
	Delete(DeleteConferenceRequest) DeleteConferenceResponse
}

// DeleteConferenceRequest is the request object for ConferenceService.Delete.
type DeleteConferenceRequest struct {
	ID uint32
}

// DeleteConferenceResponse is the response object for ConferenceService.Delete.
type DeleteConferenceResponse struct {
}

// CreateConferenceRequest is the request object for ConferenceService.Create.
type CreateConferenceRequest struct {
	Conference Conference
}

// CreateConferenceResponse is the response object for ConferenceService.Create.
type CreateConferenceResponse struct {
	Conference Conference
}

// GetConferenceRequest is the request object for ConferenceService.Get.
type GetConferenceRequest struct {
	ID uint32
}

// GetConferenceBySlugRequest is the request object for ConferenceService.GetBySlug.
type GetConferenceBySlugRequest struct {
	Slug string
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
