package main

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/gopheracademy/manager/log"
)

type conferenceService struct {
	tracer         opentracing.Tracer
	metricsFactory metrics.Factory
	logger         log.Factory
}

func newconferenceService(tracer opentracing.Tracer, metricsFactory metrics.Factory, logger log.Factory) *conferenceService {
	cs := &conferenceService{
		tracer:         tracer,
		metricsFactory: metricsFactory,
		logger:         logger,
	}
	return cs
}

func (c conferenceService) List(ctx context.Context, r ListConferenceRequest) (*ListConferenceResponse, error) {
	c.logger.For(ctx).Info("conferenceService.List")
	resp := &ListConferenceResponse{}
	return resp, nil
}

func (c conferenceService) Create(ctx context.Context, r CreateConferenceRequest) (*CreateConferenceResponse, error) {
	c.logger.For(ctx).Info("conferenceService.Create")
	resp := &CreateConferenceResponse{}
	return resp, nil
}

func (c conferenceService) Delete(ctx context.Context, r DeleteConferenceRequest) (*DeleteConferenceResponse, error) {
	c.logger.For(ctx).Info("conferenceService.Delete")
	resp := &DeleteConferenceResponse{}
	return resp, nil
}
func (c conferenceService) Get(ctx context.Context, r GetConferenceRequest) (*GetConferenceResponse, error) {
	c.logger.For(ctx).Info("conferenceService.Get")
	resp := &GetConferenceResponse{
		Conference: Conference{
			Name: "Gophercon",
			Events: []Event{{
				Name: "GopherCon 2021",
				Slug: "2021",
			},
			},
		},
	}
	return resp, nil
}

func (c conferenceService) GetBySlug(ctx context.Context, r GetConferenceBySlugRequest) (*GetConferenceResponse, error) {
	c.logger.For(ctx).Info("conferenceService.GetBySlug")
	resp := &GetConferenceResponse{
		Conference: Conference{
			Name: "Gophercon",
		},
	}
	return resp, nil
}
