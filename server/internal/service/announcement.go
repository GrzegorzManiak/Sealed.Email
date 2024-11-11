package service

import (
	"context"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func AnnounceService(ctx context.Context, client *clientv3.Client, serviceAnnouncement Announcement, value string) (clientv3.LeaseID, error) {
	if err := EnsureEtcdConnection(serviceAnnouncement.Service); err != nil {
		return 0, fmt.Errorf("failed to ensure etcd connection: %w", err)
	}

	key := serviceAnnouncement.BuildID()
	logger := helpers.GetLogger()
	logger.Printf("Registering service %s", key)

	lease, err := client.Grant(ctx, serviceAnnouncement.Service.TTL)
	if err != nil {
		return 0, fmt.Errorf("failed to create lease for service %s: %w", key, err)
	}

	_, err = client.Put(ctx, key, value, clientv3.WithLease(lease.ID))
	if err != nil {
		return 0, fmt.Errorf("failed to register service %s: %w", key, err)
	}

	logger.Printf("Service %s registered successfully", key)
	return lease.ID, nil
}

func KeepServiceAnnouncementAlive(ctx context.Context, serviceAnnouncement Announcement, unique bool) error {
	logger := helpers.GetLogger()
	marshaledService, err := serviceAnnouncement.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal service announcement: %w", err)
	}

	leaseID, err := AnnounceService(ctx, GetEtcdClient(), serviceAnnouncement, marshaledService)
	if err != nil {
		return fmt.Errorf("failed to register service %s: %w", serviceAnnouncement.Service.Prefix, err)
	}

	go func() {
		lease, err := GetEtcdClient().KeepAlive(ctx, leaseID)
		if err != nil {
			logger.Printf("failed to start KeepAlive for service %s: %v", serviceAnnouncement.Service.Prefix, err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				logger.Println("Context cancelled, stopping KeepAlive for service.")
				return

			case resp, ok := <-lease:
				if !ok {
					logger.Println("Failed to get response from KeepAlive channel, retrying.")
					lease, err = GetEtcdClient().KeepAlive(ctx, leaseID)
					if err != nil {
						logger.Printf("failed to start KeepAlive for service %s: %v", serviceAnnouncement.Service.Prefix, err)
					}
				}

				// FYI: KeepAlive handles the TTL, we just check if it failed (And restart it)
				sleepFor := (time.Duration(resp.TTL) * time.Second) / 3
				time.Sleep(sleepFor)
			}
		}
	}()

	return nil
}
