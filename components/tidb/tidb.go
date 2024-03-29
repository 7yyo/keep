package tidb

import (
	"context"
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	components "keep/components/pd"
	"keep/promp"
	"strings"
)

type TiDB struct {
	Version        string `json:"version"`
	GitHash        string `json:"git_hash"`
	Host           string
	StatusPort     int    `json:"status_port"`
	DeployPath     string `json:"deploy_path"`
	StartTimestamp int    `json:"start_timestamp"`
}

type Runner struct {
	Etcd        *clientv3.Client
	Pd          *components.PlacementDriver
	TidbCluster []*TiDB
}

var tidbOption = []string{
	". SCHEMA",
	". VARIABLES",
}

func (r *Runner) Run() error {
	p := promp.Select(tidbOption, "tidb", 20)
	i, _, err := p.Run()
	if err != nil {
		return err
	}
	switch i {
	case 0:
		return r.displayTiDBSchema()
	case 1:
		vars()
		return nil
	default:
		return nil
	}
}

func NewTiDBCluster(etcd *clientv3.Client) []*TiDB {
	r, err := etcd.Get(context.TODO(), "/topology/tidb/", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	tidbCluster := make([]*TiDB, 0)
	var tidb TiDB
	for _, v := range r.Kvs {
		if string(v.Key[len(v.Key)-4:]) == "info" {
			err := json.Unmarshal(v.Value, &tidb)
			if err != nil {
				return nil
			}
			tidb.Host = strings.Split(string(v.Key), "/")[3]
			tidbCluster = append(tidbCluster, &tidb)
		}
	}
	return tidbCluster
}
