package mapping

import (
	"context"
	"github.com/oligzeev/host-manager/internal/domain"
	v1 "github.com/openshift/api/route/v1"
	route "github.com/openshift/client-go/route/clientset/versioned"
	routeInformers "github.com/openshift/client-go/route/informers/externalversions"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kuberest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
	"sync"
)

type OpenshiftMappingService struct {
	hosts map[string]domain.Mapping

	informerStop chan struct{}
	client       *route.Clientset
	mutex        sync.RWMutex
	namespace    string
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
	client, err := route.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	routes, err := client.RouteV1().Routes(cfg.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	hosts := make(map[string]domain.Mapping)
	for _, route := range routes.Items {
		key := strings.ToLower(route.Name)
		hosts[key] = domain.Mapping{Id: key, Host: route.Spec.Host}
	}
	return &OpenshiftMappingService{hosts: hosts,
		client:       client,
		informerStop: make(chan struct{}),
		mutex:        sync.RWMutex{},
		namespace:    cfg.Namespace,
	}, nil
}

func (k *OpenshiftMappingService) StartInformer() {
	informerFactory := routeInformers.NewSharedInformerFactoryWithOptions(k.client, 0, /*30 * time.Second*/
		routeInformers.WithNamespace(k.namespace))
	informer := informerFactory.Route().V1().Routes().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    k.addEvent,
		DeleteFunc: k.deleteEvent,
		UpdateFunc: k.updateEvent,
	})
	informerFactory.Start(k.informerStop)
}

func (k *OpenshiftMappingService) addEvent(obj interface{}) {
	routeObj := obj.(*v1.Route)
	name := strings.ToLower(routeObj.Name)
	host := routeObj.Spec.Host

	k.mutex.Lock()
	mapping, ok := k.hosts[name]
	if !ok || (ok && mapping.Host != host) {
		k.hosts[name] = domain.Mapping{Id: name, Host: host}
	}
	k.mutex.Unlock()
}

func (k *OpenshiftMappingService) deleteEvent(obj interface{}) {
	routeObj := obj.(*v1.Route)
	name := strings.ToLower(routeObj.Name)

	k.mutex.Lock()
	_, ok := k.hosts[name]
	if ok {
		delete(k.hosts, name)
	}
	k.mutex.Unlock()
}

func (k *OpenshiftMappingService) updateEvent(oldObj, newObj interface{}) {
	oldRouteObj := oldObj.(*v1.Route)
	oldHost := oldRouteObj.Spec.Host
	newRouteObj := newObj.(*v1.Route)
	newName := newRouteObj.Name
	newHost := newRouteObj.Spec.Host

	if oldHost != newHost {
		k.mutex.Lock()
		mapping, ok := k.hosts[newName]
		if ok && mapping.Host != newHost {
			k.hosts[newName] = domain.Mapping{Id: newName, Host: newHost}
		}
		k.mutex.Unlock()
	}
}

func (k *OpenshiftMappingService) StopInformer() {
	const op = "OpenshiftMappingService.StopInformer"

	close(k.informerStop)
	log.Tracef("%s: finished", op)
}

// Returns all host mappings
func (k *OpenshiftMappingService) GetAll(_ context.Context, result *[]domain.Mapping) error {
	const op = "OpenshiftMappingService.GetAll"
	log.Tracef("%s", op)

	k.mutex.RLock()
	values := make([]domain.Mapping, 0, len(k.hosts))
	for _, v := range k.hosts {
		values = append(values, v)
	}
	k.mutex.RUnlock()

	*result = values
	return nil
}

// Returns host mapping by id
// Id have to be without prefix
func (k *OpenshiftMappingService) GetById(_ context.Context, id string, result *domain.Mapping) error {
	const op = "OpenshiftMappingService.GetById"
	key := strings.ToLower(id)
	log.Tracef("%s: %s", op, key)

	k.mutex.RLock()
	value, ok := k.hosts[key]
	k.mutex.RUnlock()

	if !ok {
		return domain.E(op, domain.ErrNotFound)
	}
	*result = value
	return nil
}
