package mapping

import (
	"context"
	"github.com/oligzeev/host-manager/internal/domain"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricMappingService struct {
	service                domain.MappingService
	getAllCounter          prometheus.Counter
	getByIdCounter         prometheus.Counter
	getByIdNotFoundCounter prometheus.Counter
}

func NewMetricMappingService(service domain.MappingService) *MetricMappingService {
	getAllCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "mapping_get_all_metric",
		Help: "Get all mappings counter",
	})
	getByIdCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "mapping_get_by_id_metric",
		Help: "Get mapping by id counter",
	})
	getByIdNotFoundCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "mapping_get_by_id_notfound_metric",
		Help: "Get mapping by id failed with not found error counter",
	})
	return &MetricMappingService{
		service:                service,
		getAllCounter:          getAllCounter,
		getByIdCounter:         getByIdCounter,
		getByIdNotFoundCounter: getByIdNotFoundCounter,
	}
}

func (s MetricMappingService) GetAll(ctx context.Context, result *[]domain.Mapping) error {
	s.getAllCounter.Inc()
	return s.service.GetAll(ctx, result)
}

func (s MetricMappingService) GetById(ctx context.Context, id string, result *domain.Mapping) error {
	s.getByIdCounter.Inc()
	err := s.service.GetById(ctx, id, result)
	if domain.ECode(err) == domain.ErrNotFound {
		s.getByIdNotFoundCounter.Inc()
	}
	return err
}
