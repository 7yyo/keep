package components

import (
	"go.etcd.io/etcd/clientv3"
	"time"
)

func NewEtcd(endpoints []string) (*clientv3.Client, error) {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	return cli, nil
}
