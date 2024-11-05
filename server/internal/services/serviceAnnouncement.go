package services

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func AnnounceService(ctx context.Context, client *clientv3.Client, serviceAnnouncement ServiceAnnouncement, value string) (clientv3.LeaseID, error) {
	EnsureEtcdConnection(serviceAnnouncement.Service)

	key := serviceAnnouncement.BuildID()
	logger := helpers.GetLogger()
	logger.Printf("Registering service %s", key)
	lease, err := client.Grant(ctx, serviceAnnouncement.Service.TTL)
	if err != nil {
		logger.Printf("failed to create lease for service %s: %v", key, err)
		return 0, err
	}

	_, err = client.Put(ctx, key, value, clientv3.WithLease(lease.ID))
	if err != nil {
		log.Printf("failed to register service %s: %v", key, err)
		return 0, err
	}

	return lease.ID, nil
}

func KeepServiceAnnouncementAlive(ctx context.Context, serviceAnnouncement ServiceAnnouncement, unique bool) {
	marshaledService, err := serviceAnnouncement.Marshal()
	logger := helpers.GetLogger()
	if err != nil {
		logger.Fatalf("failed to marshal service announcement: %v", err)
	}

	leaseID, err := AnnounceService(ctx, GetEtcdClient(), serviceAnnouncement, marshaledService)
	if err != nil {
		logger.Fatalf("failed to register service %s: %v", serviceAnnouncement.Service.Prefix, err)
		return
	}

	go func() {
		respChan, err := GetEtcdClient().KeepAlive(ctx, leaseID)
		if err != nil {
			logger.Printf("failed to start KeepAlive for service %s: %v", serviceAnnouncement.Service.Prefix, err)
			return
		}

		for {
			select {
			case resp, ok := <-respChan:
				if !ok {
					logger.Println("KeepAlive channel closed, stopping lease renewals.")
					KeepServiceAnnouncementAlive(ctx, serviceAnnouncement, unique)
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
