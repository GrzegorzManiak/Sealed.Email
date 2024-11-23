package service

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"time"
)

func RegisterServiceWithLease(ctx context.Context, client *clientv3.Client, key string, value string, ttl int64) (clientv3.LeaseID, error) {
	lease, err := client.Grant(ctx, ttl)
	if err != nil {
		return 0, fmt.Errorf("failed to create lease for key %s: %w", key, err)
	}

	_, err = client.Put(ctx, key, value, clientv3.WithLease(lease.ID))
	if err != nil {
		return 0, fmt.Errorf("failed to register key %s: %w", key, err)
	}

	zap.L().Info("Key registered", zap.String("key", key), zap.String("value", value))
	return lease.ID, nil
}

func AnnounceService(ctx context.Context, client *clientv3.Client, serviceAnnouncement Announcement, value string) (clientv3.LeaseID, error) {
	if err := EnsureEtcdConnection(serviceAnnouncement.Service); err != nil {
		return 0, fmt.Errorf("failed to ensure etcd connection: %w", err)
	}

	key := serviceAnnouncement.BuildID()
	ttl := serviceAnnouncement.Service.TTL

	return RegisterServiceWithLease(ctx, client, key, value, ttl)
}

func KeepServiceAnnouncementAlive(ctx context.Context, serviceAnnouncement Announcement, unique bool) error {
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
			zap.L().Error("failed to start KeepAlive for service", zap.String("service", serviceAnnouncement.Service.Prefix), zap.Error(err))
			return
		}

		for {
			select {
			case <-ctx.Done():
				zap.L().Info("Context done, stopping KeepAlive for service", zap.String("service", serviceAnnouncement.Service.Prefix))
				return

			case resp, ok := <-lease:
				if !ok {
					zap.L().Error("KeepAlive channel closed, Attempting to reconnect KeepAlive for service", zap.String("service", serviceAnnouncement.Service.Prefix))
					err = EnsureEtcdConnection(serviceAnnouncement.Service)
					if err != nil {
						zap.L().Panic("failed to ensure etcd connection", zap.Error(err))
						return
					}

					client := GetEtcdClient()
					if client == nil {
						zap.L().Panic("failed to get etcd client")
						return
					}

					leaseID, err = AnnounceService(ctx, client, serviceAnnouncement, marshaledService)
					if err != nil {
						zap.L().Panic("failed to restart KeepAlive for service", zap.String("service", serviceAnnouncement.Service.Prefix), zap.Error(err))
					}

					lease, err = client.KeepAlive(ctx, leaseID)
					if err != nil {
						zap.L().Error("failed to restart KeepAlive for service", zap.String("service", serviceAnnouncement.Service.Prefix), zap.Error(err))
					}
				}

				// FYI: KeepAlive handles the TTL, we just check if it failed (And restart it)
				timeout := serviceAnnouncement.Service.TTL
				if resp != nil {
					// -- We want to use the TTL from the response, but sometimes we dont
					//    get a response.
					timeout = resp.TTL
				}
				sleepFor := (time.Duration(timeout) * time.Second) / 3
				time.Sleep(sleepFor)
			}
		}
	}()

	return nil
}
