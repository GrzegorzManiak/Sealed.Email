package services

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
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

func LeaseService(ctx context.Context, client *clientv3.Client, serviceAnnouncement ServiceAnnouncement, value string) (clientv3.LeaseID, error) {
	key := serviceAnnouncement.BuildID()
	logger := helpers.GetLogger()
	logger.Printf("Registering service %s", key)
	lease, err := client.Grant(ctx, serviceAnnouncement.Service.TTL)
	if err != nil {
		logger.Printf("failed to create lease for service %s: %v", key, err)
		return 0, err
	}

	client = EnsureEtcdConnection(serviceAnnouncement.Service, client)

	_, err = client.Put(ctx, key, value, clientv3.WithLease(lease.ID))
	if err != nil {
		log.Printf("failed to register service %s: %v", key, err)
		return 0, err
	}

	return lease.ID, nil
}

func KeepLeaseAlive(ctx context.Context, client *clientv3.Client, serviceAnnouncement ServiceAnnouncement, unique bool) {
	marshaledService, err := serviceAnnouncement.Marshal()
	logger := helpers.GetLogger()
	if err != nil {
		logger.Fatalf("failed to marshal service announcement: %v", err)
	}

	leaseID, err := LeaseService(ctx, client, serviceAnnouncement, marshaledService)
	if err != nil {
		logger.Fatalf("failed to register service %s: %v", serviceAnnouncement.Service.Prefix, err)
		return
	}

	go func() {
		respChan, err := client.KeepAlive(ctx, leaseID)
		if err != nil {
			logger.Printf("failed to start KeepAlive for service %s: %v", serviceAnnouncement.Service.Prefix, err)
			return
		}

		for {
			select {
			case resp, ok := <-respChan:
				if !ok {
					logger.Println("KeepAlive channel closed, stopping lease renewals.")
					KeepLeaseAlive(ctx, client, serviceAnnouncement, unique)
					return
				}

				sleepFor := (time.Duration(resp.TTL) * time.Second) / 3
				time.Sleep(sleepFor)

			case <-ctx.Done():
				logger.Println("KeepLeaseAlive context canceled, exiting.")
				return
			}
		}
	}()
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
