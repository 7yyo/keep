package components

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"time"
)

func NewEtcd(endpoints []string) *clientv3.Client {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
		LogConfig: &zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.ErrorLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		},
	}
	etcd, err := clientv3.New(cfg)
	if err != nil {
		fmt.Println(err)
	}
	return etcd
}
