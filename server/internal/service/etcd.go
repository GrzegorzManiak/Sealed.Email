package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type EtcdService struct {
	client *clientv3.Client
	lock   *sync.RWMutex
	config *structs.ServiceConfig
	ans    *Announcement
	ctx    context.Context
	lease  clientv3.LeaseID
}

func NewEtcdService(ctx context.Context, config *structs.ServiceConfig, announce *Announcement) (*EtcdService, error) {
	service := &EtcdService{
		lock:   &sync.RWMutex{},
		config: config,
		ans:    announce,
		ctx:    ctx,
	}

	if err := service.keepAlive(); err != nil {
		return nil, fmt.Errorf("failed to keep alive (Issue is probably not related to etcd or this etcd service): %w", err)
	}

	zap.L().Info("Service announcement registered", zap.String("service", service.ans.String()))

	return service, nil
}

func (e *EtcdService) GetClient() (*clientv3.Client, error) {
	e.lock.RLock()
	defer e.lock.RUnlock()

	if e.client == nil {
		return nil, errors.New("etcd client is not initialized")
	}

	return e.client, nil
}

func (e *EtcdService) instantiateClient() error {
	newClient, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Etcd.Endpoints,
		DialTimeout: 2 * time.Second,
		Username:    e.config.Username,
		Password:    e.config.Password,
	})
	if err != nil {
		return fmt.Errorf("failed to instantiate etcd client: %w", err)
	}

	e.lock.Lock()
	defer e.lock.Unlock()
	e.client = newClient

	return nil
}

func (e *EtcdService) destroyClient() error {
	if e.client == nil {
		return errors.New("etcd client is already nil or uninitialized")
	}

	if err := e.client.Close(); err != nil {
		return fmt.Errorf("failed to close etcd client: %w", err)
	}

	return nil
}

func (e *EtcdService) CheckClientConnection() error {
	if e.client == nil {
		return errors.New("etcd client is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for _, endpoint := range e.client.Endpoints() {
		if _, err := e.client.Status(ctx, endpoint); err != nil {
			return fmt.Errorf("etcd client connection failed for endpoint %s: %w", endpoint, err)
		}
	}

	return nil
}

func (e *EtcdService) EnsureConnection() error {
	if e.client != nil {
		e.lock.Lock()
		defer e.lock.Unlock()

		if err := e.CheckClientConnection(); err == nil {
			return nil
		} else {
			zap.L().Warn("etcd client connection failed, destroying client")

			if err := e.destroyClient(); err != nil {
				zap.L().Warn("failed to destroy etcd client", zap.Error(err))
			}

			e.client = nil
		}
	}

	zap.L().Info("no etcd client, instantiating new client")

	if err := e.instantiateClient(); err != nil {
		return fmt.Errorf("failed to ensure etcd connection: %w", err)
	}

	return nil
}

func (e *EtcdService) GetAllKeys(ctx context.Context) ([]*mvccpb.KeyValue, error) {
	if e.client == nil {
		return nil, errors.New("etcd client is not initialized")
	}

	resp, err := e.client.Get(ctx, Prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve keys with prefix %s: %w", Prefix, err)
	}

	return resp.Kvs, nil
}
