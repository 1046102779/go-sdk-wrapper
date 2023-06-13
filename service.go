package wrapper

import (
	"errors"
	"reflect"

	"github.com/dapr/go-sdk/service/common"
	services "github.com/dapr/go-sdk/service/grpc"
	"github.com/dapr/kit/logger"
)

// Service service handler object.
type Service struct {
	daprSrv        common.Service
	adapterDaprAPI *adapterDaprAPI
}

// NewService create a new dapr service.
func NewService(log *logger.Logger, addr string) (srv *Service, err error) {
	var daprSrv common.Service
	daprSrv, err = services.NewService(addr)
	if err != nil {
		return nil, err
	}

	return &Service{
		daprSrv:        daprSrv,
		adapterDaprAPI: newAdapterDaprAPI(log),
	}, nil
}

// RegisterRPC registers MicroLogic RPC API.
func (s *Service) RegisterRPC(router string, fn interface{}) error {
	if fn == nil {
		return errors.New("service handle func empty")
	}
	val := reflect.ValueOf(fn)

	return s.daprSrv.AddServiceInvocationHandler(router, s.adapterDaprAPI.RPC(val))
}

// Subscription simplfy pubsub.sub.
type Subscription struct {
	// PubsubName is name of the pub/sub this message came from
	PubsubName string `json:"pubsubname"`
	// Topic is the name of the topic
	Topic string `json:"topic"`
	// Metadata is the subscription metadata
	Metadata map[string]string `json:"metadata,omitempty"`
}

// RegisterSubscribe registers subscribe msg.
func (s *Service) RegisterSubscribe(sub *Subscription, fn interface{}) error {
	val := reflect.ValueOf(fn)

	return s.daprSrv.AddTopicEventHandler(&common.Subscription{
		PubsubName: sub.PubsubName,
		Topic:      sub.Topic,
		Metadata:   sub.Metadata,
	}, s.adapterDaprAPI.Subscribe(val))
}

// RegisterInput registers inputbindings event.
func (s *Service) RegisterInput(router string, fn interface{}) error {
	val := reflect.ValueOf(fn)

	return s.daprSrv.AddBindingInvocationHandler(router,
		s.adapterDaprAPI.InputBinding(val))
}

// Start starts server.
func (s *Service) Start() error {
	return s.daprSrv.Start()
}
