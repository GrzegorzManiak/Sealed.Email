package service

import (
	"context"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type Pools struct {
	etcdService    *EtcdService
	poolsMutex     sync.RWMutex
	pools          map[string]*Pool
	ctx            context.Context
	securityPolicy grpc.DialOption
}

func NewPools(ctx context.Context, etcdService *EtcdService, certs structs.ServiceCertificates) (*Pools, error) {
	pools := &Pools{
		etcdService:    etcdService,
		pools:          make(map[string]*Pool),
		ctx:            ctx,
		securityPolicy: GetTransportSecurityPolicy(certs),
	}
	pools.keepPoolsFresh()
	return pools, nil
}

func (p *Pools) refreshPools() error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	keyValues, err := p.etcdService.GetAllKeys(ctx)
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	p.poolsMutex.Lock()
	defer p.poolsMutex.Unlock()
	newPools := make(map[string]*Pool)

	for _, keyValue := range keyValues {
		announcement, err := UnmarshalServiceAnnouncement(keyValue.Value)
		if err != nil {
			zap.L().Warn("failed to unmarshal service announcement", zap.Error(err))
			continue
		}

		pool, exists := newPools[announcement.Service.Prefix]
		if !exists {
			pool = &Pool{
				Mutex:    sync.RWMutex{},
				Pool:     make(map[string]*GrpcConnection),
				Security: p.securityPolicy,
				Keys:     []string{},
			}
			newPools[announcement.Service.Prefix] = pool
		}

		pool.Mutex.Lock()
		emptyConnection := &GrpcConnection{Announcement: announcement}
		pool.Pool[announcement.Id] = emptyConnection
		pool.Keys = append(pool.Keys, announcement.Id)
		pool.Mutex.Unlock()
	}

	p.pools = newPools
	return nil
}

func (p *Pools) keepPoolsFresh() {
	go func() {
		zap.L().Info("Starting pool refresh")
		time.Sleep(2 * time.Second)
		for {
			select {
			case <-p.ctx.Done():
				zap.L().Warn("Context done, stopping pool refresh", zap.Error(p.ctx.Err()))
				return

			default:
				if err := p.refreshPools(); err != nil {
					zap.L().Error("failed to refresh pools", zap.Error(err))
				}
				time.Sleep(time.Duration(config.Etcd.ConnectionPool.RefreshInterval) * time.Second)
			}
		}
	}()
}

func (p *Pools) GetPool(servicePrefix string) (*Pool, error) {
	p.poolsMutex.RLock()
	pool, exists := p.pools[servicePrefix]
	p.poolsMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("pool for service %s not found", servicePrefix)
	}

	return pool, nil
}
