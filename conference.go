package main

import (
	"context"
	"log"
)

type conferenceService struct{}

func (conferenceService) List(ctx context.Context, r ListConferenceRequest) (*ListConferenceResponse, error) {

	log.Println("list")
	resp := &ListConferenceResponse{}
	return resp, nil
}

func (conferenceService) Get(ctx context.Context, r GetConferenceRequest) (*GetConferenceResponse, error) {
	log.Println("get")
	resp := &GetConferenceResponse{
		Conference: Conference{
			Name: "Gophercon",
		},
	}
	return resp, nil
}
