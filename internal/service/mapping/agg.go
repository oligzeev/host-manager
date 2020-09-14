package mapping

import (
	"context"
	"github.com/oligzeev/host-manager/internal/domain"
	log "github.com/sirupsen/logrus"
)

type AggMappingService struct {
	envMappingService       domain.MappingService
	openshiftMappingService domain.MappingService
}

func NewAggMappingService(
	envMappingService domain.MappingService,
	openshiftMappingService domain.MappingService,
) *AggMappingService {
	return &AggMappingService{envMappingService: envMappingService, openshiftMappingService: openshiftMappingService}
}

func (s AggMappingService) GetAll(ctx context.Context, result *[]domain.Mapping) error {
	const op = "AggMappingService.GetAll"

	log.Tracef("%s", op)
	var err error
	var envHosts, openshiftHosts []domain.Mapping
	err = s.envMappingService.GetAll(ctx, &envHosts)
	if err != nil {
		return err
	}
	if s.openshiftMappingService != nil {
		err = s.openshiftMappingService.GetAll(ctx, &openshiftHosts)
		if err != nil {
			return err
		}
	}
	var hosts []domain.Mapping
	if len(envHosts) > 0 {
		hosts = append(hosts, envHosts...)
	}
	if s.openshiftMappingService != nil && len(openshiftHosts) > 0 {
		hosts = append(hosts, openshiftHosts...)
	}
	*result = hosts
	return nil
}

func (s AggMappingService) GetById(ctx context.Context, id string, result *domain.Mapping) error {
	const op = "AggMappingService.GetById"

	log.Tracef("%s", op)
	var err error
	err = s.envMappingService.GetById(ctx, id, result)
	if err != nil && domain.ECode(err) == domain.ErrNotFound && s.openshiftMappingService != nil {
		err = s.openshiftMappingService.GetById(ctx, id, result)
		if err == nil {
			return nil
		}
	}
	return err
}
