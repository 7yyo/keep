package components

import (
	"go.etcd.io/etcd/clientv3"
	"time"
)

func NewEtcd(endpoints []string) *clientv3.Client {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}
	etcd, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}
	return etcd
}
