/*
Package service provides very little - it has a service.Registry which can act
as a virtual service (aggregate of endpoints).
*/
package service

import (
	"fmt"
	"sync"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
)

type Registry struct {
	services map[string]api.Service
	sync.Mutex
}

var defaultRegistry = &Registry{}

func (r *Registry) Register(name string, s api.Service) {
	r.Lock()
	defer r.Unlock()
	if r.services == nil {
		r.services = make(map[string]api.Service)
	}

	if _, ok := r.services[name]; ok {
		panic(fmt.Sprintf("service already registered :%s", name))
	}
	r.services[name] = s
}

func (r *Registry) New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	endpoint := make(api.Endpoint)

	for serviceName, s := range r.services {
		for name, reqBuilder := range s.New(p, opts...) {
			endpoint[serviceName+":"+name] = reqBuilder
		}
	}
	return endpoint
}

func Register(name string, s api.Service) {
	defaultRegistry.Register(name, s)
}

func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	return defaultRegistry.New(p, opts...)
}
