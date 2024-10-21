package internal

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
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
		log.Printf("failed to close etcd client: %v", err)
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
	if existingClient != nil {
		err := CheckClientConnection(existingClient)
		if err == nil {
			log.Printf("etcd connection successful")
			return existingClient
		}
		log.Printf("etcd connection failed: %v", err)
		DestroyEtcdClient(existingClient)
	}

	client := GetEtcdClient(service)
	log.Printf("etcd connection successful")
	return client
}

func LeaseService(ctx context.Context, client *clientv3.Client, service structs.ServiceConfig, value string) (clientv3.LeaseID, error) {
	log.Printf("Registering service %s", service.Prefix)
	lease, err := client.Grant(ctx, service.TTL)
	if err != nil {
		log.Printf("failed to create lease for service %s: %v", service.Prefix, err)
		return 0, err
	}

	client = EnsureEtcdConnection(service, client)

	_, err = client.Put(ctx, service.Prefix, value, clientv3.WithLease(lease.ID))
	if err != nil {
		log.Printf("failed to register service %s: %v", service.Prefix, err)
		return 0, err
	}

	return lease.ID, nil
}

func KeepLeaseAlive(ctx context.Context, client *clientv3.Client, service structs.ServiceConfig, value string) {
	leaseID, err := LeaseService(ctx, client, service, value)
	if err != nil {
		log.Fatalf("failed to register service %s: %v", service.Prefix, err)
		return
	}

	go func() {
		respChan, err := client.KeepAlive(ctx, leaseID)
		if err != nil {
			log.Printf("failed to start KeepAlive for service %s: %v", service.Prefix, err)
			return
		}

		for {
			select {
			case resp, ok := <-respChan:
				if !ok {
					log.Println("KeepAlive channel closed, stopping lease renewals.")
					KeepLeaseAlive(ctx, client, service, value)
					return
				}

				sleepFor := time.Duration(resp.TTL) * time.Second / 3
				time.Sleep(sleepFor)

			case <-ctx.Done():
				log.Println("KeepLeaseAlive context canceled, exiting.")
				return
			}
		}
	}()
}
