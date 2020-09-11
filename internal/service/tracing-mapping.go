package service

import (
	"context"
	"github.com/oligzeev/host-manager/internal/domain"
	"github.com/opentracing/opentracing-go"
)

type TracingMappingService struct {
	service domain.MappingService
}

func NewTracingMappingService(service domain.MappingService) *TracingMappingService {
	return &TracingMappingService{service: service}
}

func (s TracingMappingService) GetAll(ctx context.Context, result *[]domain.Mapping) error {
	const op = "MappingService.GetAll"
	span, spanCtx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()
	return s.service.GetAll(spanCtx, result)
}

func (s TracingMappingService) GetById(ctx context.Context, id string, result *domain.Mapping) error {
	const op = "MappingService.GetById"
	span, spanCtx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()
	return s.service.GetById(spanCtx, id, result)
}
