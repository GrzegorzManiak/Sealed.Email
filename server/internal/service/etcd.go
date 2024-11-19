package service

import (
	"context"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"sync"
	"time"
)

var clientLock = &sync.RWMutex{}
var client = &clientv3.Client{}

func GetEtcdClient() *clientv3.Client {
	clientLock.RLock()
	defer clientLock.RUnlock()
	if client == nil {
		zap.L().Panic("etcd client is not initialized")
	}
	return client
}

func InstantiateEtcdClient(service structs.ServiceConfig) error {
	newClient, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Etcd.Endpoints,
		DialTimeout: 2 * time.Second,
		Username:    service.Username,
		Password:    service.Password,
	})

	if err != nil {
		return fmt.Errorf("failed to instantiate etcd client: %w", err)
	}

	clientLock.Lock()
	defer clientLock.Unlock()
	client = newClient

	return nil
}

func DestroyEtcdClient(client *clientv3.Client) error {
	if client == nil {
		return fmt.Errorf("etcd client is already nil or uninitialized")
	}

	if err := client.Close(); err != nil {
		return fmt.Errorf("failed to close etcd client: %w", err)
	}

	return nil
}

func CheckClientConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for _, endpoint := range client.Endpoints() {
		if _, err := client.Status(ctx, endpoint); err != nil {
			return fmt.Errorf("etcd client connection failed for endpoint %s: %w", endpoint, err)
		}
	}

	return nil
}

func EnsureEtcdConnection(service structs.ServiceConfig) error {
	clientLock.Lock()
	defer clientLock.Unlock()

	if client != nil {
		if err := CheckClientConnection(); err == nil {
			return nil
		} else {
			zap.L().Warn("etcd client connection failed, reconnecting", zap.Error(err))
			if err := DestroyEtcdClient(client); err != nil {
				return fmt.Errorf("failed to destroy etcd client: %w", err)
			}
			client = nil
		}
	}

	if err := InstantiateEtcdClient(service); err != nil {
		zap.L().Error("failed to ensure etcd connection", zap.Error(err))
		return fmt.Errorf("failed to ensure etcd connection: %w", err)
	}

	return nil
}

func GetAllKeys(ctx context.Context, client *clientv3.Client) ([]*mvccpb.KeyValue, error) {
	resp, err := client.Get(ctx, Prefix, clientv3.WithPrefix())
	if err != nil {
		zap.L().Error("failed to retrieve keys", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve keys with prefix %s: %w", Prefix, err)
	}
	return resp.Kvs, nil
}
