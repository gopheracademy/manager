package server

import (
	"context"

	"github.com/gopheracademy/manager/backend/generated"
)

type ConferenceService struct{}

func (ConferenceService) List(ctx context.Context, r generated.ListConferenceRequest) (*generated.ListConferenceResponse, error) {
	resp := &generated.ListConferenceResponse{}
	return resp, nil
}

func (ConferenceService) Get(ctx context.Context, r generated.GetConferenceRequest) (*generated.GetConferenceResponse, error) {
	resp := &generated.GetConferenceResponse{
		Conference: generated.Conference{
			Name: "Gophercon",
		},
	}
	return resp, nil
}
