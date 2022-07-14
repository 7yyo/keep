package components

import (
	"go.etcd.io/etcd/clientv3"
	"keep/components/cdc"
	components "keep/components/pd"
	"keep/components/tidb"
)

type Server struct {
	Etcd *clientv3.Client
	*components.PlacementDriver
	TiDBCluster []*tidb.TiDB
}

type Runner interface {
	Run() error
}

func (s *Server) Run(c string) error {
	var err error
	switch c {
	case "cdc":
		r := cdc.Runner{
			Etcd: s.Etcd,
			Pd:   s.PlacementDriver,
		}
		err = r.Run()
	case "tidb":
		r := tidb.Runner{
			Etcd:        s.Etcd,
			Pd:          s.PlacementDriver,
			TidbCluster: s.TiDBCluster,
		}
		err = r.Run()
	}
	if err != nil {
		return err
	}
	return nil
}
