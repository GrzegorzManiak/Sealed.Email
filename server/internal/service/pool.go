package service

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Pool struct {
	Mutex    sync.RWMutex
	Pool     map[string]*GrpcConnection
	Security grpc.DialOption
	Keys     []string
}

type GrpcConnection struct {
	Conn         *grpc.ClientConn
	Announcement Announcement

	LastRefreshed int64
	LastChecked   int64
	Working       bool
}

func (g *GrpcConnection) initializeGrpcConnection(pool *Pool) error {
	if g.Conn != nil {
		return nil
	}

	if g.Conn != nil && !g.Working {
		zap.L().Info("Connection is not working, closing connection", zap.String("service", g.Announcement.Service.Prefix))

		err := g.Conn.Close()
		if err != nil {
			zap.L().Warn("failed to close connection!", zap.Error(err))
		}

		g.Conn = nil
	}

	if pool.Security == nil {
		return errors.New("security option not provided")
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", g.Announcement.Host, g.Announcement.Port),
		pool.Security,
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*1),
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    10 * time.Second, // Check connection every 10 seconds
			Timeout: 5 * time.Second,  // Timeout after 5 seconds of no response
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	curTime := time.Now().Unix()
	g.Conn = conn
	g.LastRefreshed = curTime
	g.LastChecked = curTime
	g.Working = true

	return nil
}

func (p *Pool) GetConnection() (*GrpcConnection, error) {
	p.Mutex.RLock()
	defer p.Mutex.RUnlock()

	if len(p.Pool) == 0 {
		return nil, errors.New("no connections available")
	}

	for range 3 {
		randKey := p.Keys[rand.Intn(len(p.Keys))]
		conn := p.Pool[randKey]

		if err := conn.initializeGrpcConnection(p); err != nil {
			zap.L().Warn("Failed to initialize connection", zap.String("key", randKey), zap.Error(err))

			continue
		}

		return conn, nil
	}

	return nil, errors.New("no connections available (tried 3 times)")
}
