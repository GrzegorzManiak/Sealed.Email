package service

import (
	"context"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

var clientLock = &sync.RWMutex{}
var client = &clientv3.Client{}

func GetEtcdClient() *clientv3.Client {
	clientLock.RLock()
	defer clientLock.RUnlock()
	if client == nil {
		panic("etcd client is not initialized")
	}
	return client
}

func InstantiateEtcdClient(service structs.ServiceConfig) error {
	if client != nil {
		helpers.GetLogger().Printf("[WARN] etcd client already initialized")
		return nil
	}

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

func EnsureEtcdConnection(service structs.ServiceConfig) {
	clientLock.Lock()
	defer clientLock.Unlock()
	logger := helpers.GetLogger()

	if client != nil {
		if err := CheckClientConnection(); err == nil {
			logger.Println("etcd connection successful")
			return
		} else {
			logger.Printf("etcd connection lost, attempting to reconnect: %v", err)
			if err := DestroyEtcdClient(client); err != nil {
				logger.Printf("failed to destroy etcd client: %v", err)
			}
			client = nil
		}
	}

	logger.Println("connecting to etcd")
	if err := InstantiateEtcdClient(service); err != nil {
		logger.Printf("failed to reconnect to etcd: %v", err)
	} else {
		logger.Println("reconnected to etcd successfully")
	}
}

func GetAllKeys(ctx context.Context, client *clientv3.Client) ([]*mvccpb.KeyValue, error) {
	resp, err := client.Get(ctx, Prefix, clientv3.WithPrefix())
	if err != nil {
		helpers.GetLogger().Printf("failed to get keys: %v", err)
		return nil, err
	}

	return resp.Kvs, nil
}
