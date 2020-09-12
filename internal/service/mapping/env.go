package mapping

import (
	"context"
	"github.com/oligzeev/host-manager/internal/domain"
	log "github.com/sirupsen/logrus"
	"strings"
)

type EnvMappingService struct {
	hosts   []domain.Mapping
	hostMap map[string]domain.Mapping
}

// Instantiate host mapping service with specific prefix
// There's no separator for prefix so you have to specify one if required (e.g. "someprefix_")
// Prefix and all keys will be fixed with upper case
func NewEnvMappingService(cfg domain.MappingConfig, envs []string) *EnvMappingService {
	prefix := strings.ToLower(cfg.Prefix)
	hostMap := make(map[string]domain.Mapping)
	var hosts []domain.Mapping
	for _, env := range envs {
		pair := strings.SplitN(env, "=", 2)
		key := strings.ToLower(pair[0])
		value := pair[1]
		if strings.HasPrefix(key, prefix) {
			id := strings.TrimPrefix(key, prefix)
			mapping := domain.Mapping{Id: id, Host: value}
			hostMap[id] = mapping
			hosts = append(hosts, mapping)
		}
	}
	return &EnvMappingService{hostMap: hostMap, hosts: hosts}
}

// Returns all host mappings
func (k EnvMappingService) GetAll(_ context.Context, result *[]domain.Mapping) error {
	const op = "EnvMappingService.GetAll"

	log.Tracef("%s", op)
	*result = k.hosts
	return nil
}

// Returns host mapping by id
// Id have to be without prefix
func (k EnvMappingService) GetById(_ context.Context, id string, result *domain.Mapping) error {
	const op = "EnvMappingService.GetById"

	key := strings.ToLower(id)
	value, ok := k.hostMap[key]
	log.Tracef("%s: %s, %t", op, key, ok)
	if !ok {
		return domain.E(op, domain.ErrNotFound)
	}
	*result = value
	return nil
}
