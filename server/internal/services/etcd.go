package services

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func GetEtcdClient(service structs.ServiceConfig) *clientv3.Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Etcd.Endpoints,
		DialTimeout: 2 * time.Second,
		Username:    service.Username,
		Password:    service.Password,
	})

	if err != nil {
		panic(err)
	}

	return client
}

func DestroyEtcdClient(client *clientv3.Client) {
	err := client.Close()
	if err != nil {
		helpers.GetLogger().Printf("failed to close etcd client: %v", err)
	}
}

func CheckClientConnection(client *clientv3.Client) error {
	for _, endpoint := range client.Endpoints() {
		_, err := client.Status(context.Background(), endpoint)
		if err != nil {
			return err
		}
	}
	return nil
}

func EnsureEtcdConnection(service structs.ServiceConfig, existingClient *clientv3.Client) *clientv3.Client {
	logger := helpers.GetLogger()
	if existingClient != nil {
		err := CheckClientConnection(existingClient)
		if err == nil {
			logger.Printf("etcd connection successful")
			return existingClient
		}
		logger.Printf("etcd connection failed: %v", err)
		DestroyEtcdClient(existingClient)
	}

	client := GetEtcdClient(service)
	logger.Printf("etcd connection successful")
	return client
}

func GetAllLeases(ctx context.Context, client *clientv3.Client) ([]clientv3.LeaseStatus, error) {

	resp, err := client.Leases(ctx)
	if err != nil {
		return nil, err
	}

	return resp.Leases, nil
}

type KeyValueImpl struct {
	Key   string
	Value string
}

func GetAllKeys(ctx context.Context, client *clientv3.Client) ([]KeyValueImpl, error) {
	resp, err := client.Get(ctx, ServicePrefix, clientv3.WithPrefix())
	if err != nil {
		helpers.GetLogger().Printf("failed to get keys: %v", err)
		return nil, err
	}

	keys := make([]KeyValueImpl, 0)

	for _, kv := range resp.Kvs {
		keys = append(keys, KeyValueImpl{Key: string(kv.Key), Value: string(kv.Value)})
	}

	return keys, nil
}
