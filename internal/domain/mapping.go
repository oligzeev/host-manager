package domain

import "context"

type Mapping struct {
	Id   string `json:"id"`
	Host string `json:"host"`
}

type MappingService interface {
	GetAll(ctx context.Context, result *[]Mapping) error
	GetById(ctx context.Context, id string, result *Mapping) error
}
