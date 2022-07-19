package components

import (
	"go.etcd.io/etcd/clientv3"
	"keep/components/cdc"
	p "keep/components/pd"
	"keep/components/tidb"
)

type Server struct {
	Etcd *clientv3.Client
	*p.PlacementDriver
	TiDBCluster []*tidb.TiDB
}

type Runner interface {
	Run() error
}

func (s *Server) Run(c string) error {
	switch c {
	case "cdc":
		r := cdc.Runner{
			Etcd: s.Etcd,
			Pd:   s.PlacementDriver,
		}
		return r.Run()
	case "tidb":
		r := tidb.Runner{
			Etcd:        s.Etcd,
			Pd:          s.PlacementDriver,
			TidbCluster: s.TiDBCluster,
		}
		return r.Run()
	case "pd":
		r := p.Runner{
			Pd: s.PlacementDriver,
		}
		return r.Run()
	default:
		return nil
	}
}
