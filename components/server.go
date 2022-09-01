package components

import (
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"keep/components/cdc"
	p "keep/components/pd"
	"keep/components/tidb"
	"keep/util/printer"
	"time"
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
			Etcd:            s.Etcd,
			PlacementDriver: s.PlacementDriver,
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
			Pd:   s.PlacementDriver,
			Etcd: s.Etcd,
		}
		return r.Run()
	}
	return nil
}

func NewEtcd(endpoints []string) *clientv3.Client {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
		LogConfig: &zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.FatalLevel),
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
		printer.PrintError(err.Error())
	}
	return etcd
}
