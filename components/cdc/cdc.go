package cdc

import (
	"go.etcd.io/etcd/clientv3"
	components "keep/components/pd"
	"keep/promp"
	"sync"
)

type Runner struct {
	captures     []capture
	changefeedId string
	Etcd         *clientv3.Client
	Pd           *components.PlacementDriver
	sync.RWMutex
}

var cdcOption = []string{"capture", "changefeed"}

func (r *Runner) Run() error {
	p := promp.Select(cdcOption, "cdc capture list", 20)
	i, _, err := p.Run()
	if err != nil {
		return err
	}
	switch i {
	case 0:
		return r.displayCapture()
	case 1:
		return r.displayChangefeedList()
	default:
		return nil
	}
}
