package cdc

import (
	"github.com/7yyo/sunflower/prompt"
	"go.etcd.io/etcd/clientv3"
	components "keep/components/pd"
	"sync"
)

type Runner struct {
	captures        []capture
	changefeedId    string
	Etcd            *clientv3.Client
	PlacementDriver *components.PlacementDriver
	sync.RWMutex
}

func (r *Runner) Run() error {
	se := prompt.Select{
		Title: "ticdc:",
		Option: []interface{}{
			"capture",
			"changefeed",
		},
	}
	i, _, err := se.Run()
	if err != nil {
		if prompt.IsBackSpace(err) {
			return r.Run()
		}
		return err
	}
	switch i {
	case 0:
		return r.displayCaptureList()
	case 1:
		return r.displayChangefeedList()
	}
	return nil
}
