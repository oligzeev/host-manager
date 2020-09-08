package service

import (
	"context"
	"github.com/oligzeev/host-manager/internal/domain"
	log "github.com/sirupsen/logrus"
	"strings"

	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kuberest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type OpenshiftMappingService struct {
	hosts   []domain.Mapping
	hostMap map[string]domain.Mapping
}

func NewOpenshiftMappingService(cfg domain.MappingConfig) (*OpenshiftMappingService, error) {
	config, err := kuberest.InClusterConfig()
	if err != nil {
		log.Errorf("In cluster config failed: %v", err)
		config, err = clientcmd.BuildConfigFromFlags("", cfg.ConfigPath)
		if err != nil {
			return nil, err
		}
	}
	routeClient, err := routev1.NewForConfig(config)

	routes, err := routeClient.Routes(cfg.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	hostMap := make(map[string]domain.Mapping)
	var hosts []domain.Mapping
	for _, route := range routes.Items {
		key := strings.ToLower(route.Name)
		value := route.Spec.Host
		mapping := domain.Mapping{Id: key, Host: value}
		hostMap[key] = mapping
		hosts = append(hosts, mapping)
	}
	return &OpenshiftMappingService{hostMap: hostMap, hosts: hosts}, nil
}

// Returns all host mappings
func (k OpenshiftMappingService) GetAll(_ context.Context, result *[]domain.Mapping) error {
	const op = "OpenshiftMappingService.GetAll"

	log.Tracef("%s", op)
	*result = k.hosts
	return nil
}

// Returns host mapping by id
// Id have to be without prefix
func (k OpenshiftMappingService) GetById(_ context.Context, id string, result *domain.Mapping) error {
	const op = "OpenshiftMappingService.GetById"

	key := strings.ToLower(id)
	value, ok := k.hostMap[key]
	log.Tracef("%s: %s, %t", op, key, ok)
	if !ok {
		return domain.E(op, domain.ErrNotFound)
	}
	*result = value
	return nil
}
