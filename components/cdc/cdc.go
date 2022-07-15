package cdc

import (
	"github.com/manifoldco/promptui"
	"go.etcd.io/etcd/clientv3"
	components "keep/components/pd"
	"sync"
)

type Runner struct {
	captures     []capture
	changefeedId string
	Etcd         *clientv3.Client
	Pd           *components.PlacementDriver
	sync.RWMutex
}

func (r *Runner) Run() error {

	prompt := promptui.Select{
		Label: "cdc",
		Items: []string{
			"capture",
			"changefeed",
		},
	}
	_, m, err := prompt.Run()
	if err != nil {
		return err
	}

	switch m {
	case "capture":
		return r.displayCapture()
	case "changefeed":
		return r.displayChangefeed()
	}
	return nil
}
