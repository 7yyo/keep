package tidb

import (
	"context"
	"encoding/json"
	"github.com/manifoldco/promptui"
	"go.etcd.io/etcd/clientv3"
	components "keep/components/pd"
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

func (r *Runner) Run() error {

	prompt := promptui.Select{
		Label: "tidb",
		Items: []string{
			"schema",
		},
	}
	_, m, err := prompt.Run()
	if err != nil {
		return err
	}

	switch m {
	case "schema":
		err = r.displayTiDBSchema()
	}
	if err != nil {
		return err
	}
	return nil
}

func NewTiDBCluster(etcd *clientv3.Client) []*TiDB {
	r, err := etcd.Get(context.TODO(), "/topology/tidb/", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}

	tidbCluster := make([]*TiDB, 0)
	var t TiDB
	for _, v := range r.Kvs {
		if string(v.Key[len(v.Key)-4:]) == "info" {
			err := json.Unmarshal(v.Value, &t)
			if err != nil {
				return nil
			}
			t.Host = strings.Split(string(v.Key), "/")[3]
			tidbCluster = append(tidbCluster, &t)
		}
	}
	return tidbCluster
}
