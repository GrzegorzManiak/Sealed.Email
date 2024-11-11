package service

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

var clientLock = &sync.Mutex{}
var client = &clientv3.Client{}

func GetEtcdClient() *clientv3.Client {
	clientLock.Lock()
	defer clientLock.Unlock()
	if client == nil {
		panic("etcd client is not initialized")
	}
	return client
}

func InstantiateEtcdClient(service structs.ServiceConfig) {
	if client != nil {
		helpers.GetLogger().Printf("[WARN] etcd client already initialized")
	}

	newClient, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Etcd.Endpoints,
		DialTimeout: 2 * time.Second,
		Username:    service.Username,
		Password:    service.Password,
	})

	if err != nil {
		panic(err)
	}

	clientLock.Lock()
	defer clientLock.Unlock()
	client = newClient
}

func DestroyEtcdClient(client *clientv3.Client) {
	err := client.Close()
	if err != nil {
		helpers.GetLogger().Printf("failed to close etcd client: %v", err)
	}
}

func CheckClientConnection() error {
	for _, endpoint := range client.Endpoints() {
		_, err := client.Status(context.Background(), endpoint)
		if err != nil {
			return err
		}
	}
	return nil
}

func EnsureEtcdConnection(service structs.ServiceConfig) {
	logger := helpers.GetLogger()
	if client != nil {
		err := CheckClientConnection()
		if err == nil {
			logger.Printf("etcd connection successful")
			return
		}
		logger.Printf("etcd connection failed: %v", err)
		DestroyEtcdClient(client)
	}

	logger.Printf("reconnecting to etcd")
	InstantiateEtcdClient(service)
}

func GetAllLeases(ctx context.Context, client *clientv3.Client) ([]clientv3.LeaseStatus, error) {

	resp, err := client.Leases(ctx)
	if err != nil {
		return nil, err
	}

	return resp.Leases, nil
}

func GetAllKeys(ctx context.Context, client *clientv3.Client) ([]*mvccpb.KeyValue, error) {
	resp, err := client.Get(ctx, Prefix, clientv3.WithPrefix())
	if err != nil {
		helpers.GetLogger().Printf("failed to get keys: %v", err)
		return nil, err
	}

	return resp.Kvs, nil
}
