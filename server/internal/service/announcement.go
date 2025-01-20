package service

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"time"
)

func (e *EtcdService) registerLease() (clientv3.LeaseID, error) {
	client, err := e.GetClient()
	if err != nil {
		return 0, fmt.Errorf("failed to get etcd client: %w", err)
	}
	lease, err := client.Grant(e.ctx, e.ans.Service.TTL)
	if err != nil {
		return 0, fmt.Errorf("failed to register lease: %w", err)
	}
	e.lease = lease.ID
	return e.lease, nil
}

func (e *EtcdService) registerKeyValue(key string, value string) error {
	if e.lease == 0 {
		return fmt.Errorf("lease not registered")
	}

	client, err := e.GetClient()
	if err != nil {
		return fmt.Errorf("failed to get etcd client: %w", err)
	}

	_, err = client.Put(e.ctx, key, value, clientv3.WithLease(e.lease))
	if err != nil {
		return fmt.Errorf("failed to register key %s: %w", key, err)
	}

	return nil
}

func (e *EtcdService) registerService() error {
	zap.L().Info("Lease not registered, attempting to register lease for service", zap.String("service", e.ans.Service.Prefix))
	_, err := e.registerLease()
	if err != nil {
		return fmt.Errorf("failed to register lease for service: %w", err)
	}

	// -- Register the key-value pair with the lease
	marshaledService, err := e.ans.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal service announcement: %w", err)
	}

	key := e.ans.BuildID()
	if err := e.registerKeyValue(key, marshaledService); err != nil {
		return fmt.Errorf("failed to register key-value pair: %w", err)
	}

	zap.L().Info("Key registered", zap.String("key", key), zap.String("value", marshaledService))
	return nil
}

func (e *EtcdService) keepAlive() error {
	go func() {
		for {

			// -- Check if the context is done
			select {
			case <-e.ctx.Done():
				zap.L().Info("Context done, stopping KeepAlive for service", zap.String("service", e.ans.Service.Prefix))
				return

			default:
				zap.L().Info("Starting KeepAlive for service", zap.String("service", e.ans.Service.Prefix))

				// -- Ensure that the client is instantiated & connected
				if err := e.EnsureConnection(); err != nil {
					zap.L().Warn("failed to ensure etcd connection", zap.Error(err))
					e.sleep(e.ans.Service.TimeOut)
					continue
				}

				client, err := e.GetClient()
				if err != nil {
					zap.L().Warn("failed to get etcd client", zap.Error(err))
					e.sleep(e.ans.Service.TimeOut)
					continue
				}

				// -- Register the service
				err = e.registerService()
				if err != nil {
					zap.L().Warn("failed to register service", zap.String("service", e.ans.Service.Prefix), zap.Error(err))
					e.sleep(e.ans.Service.TimeOut)
					continue
				}

				// -- Send a keep alive signal
				aliveChanel, err := client.KeepAlive(e.ctx, e.lease)
				if err != nil {
					zap.L().Warn("failed to send keep alive signal for service", zap.String("service", e.ans.Service.Prefix), zap.Error(err))
					e.sleep(e.ans.Service.TimeOut)
					continue
				}

				for {
					_, ok := <-aliveChanel
					if !ok {
						zap.L().Warn("KeepAlive channel closed, Attempting to reconnect KeepAlive for service", zap.String("service", e.ans.Service.Prefix))
						break
					}
				}

				zap.L().Info("KeepAlive failed, attempting to reconnect KeepAlive for service", zap.String("service", e.ans.Service.Prefix))
				e.lease = 0
				e.sleep(e.ans.Service.TimeOut)
			}
		}
	}()

	return nil
}

func (e *EtcdService) sleep(seconds int64) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
