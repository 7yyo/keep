package cdc

import (
	"github.com/manifoldco/promptui"
	"go.etcd.io/etcd/clientv3"
	components "keep/components/pd"
)

type Runner struct {
	captures     []capture
	changefeedId string
	Etcd         *clientv3.Client
	Pd           *components.PlacementDriver
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
		if err := r.DisplayCapture(); err != nil {
			return err
		}
	case "changefeed":
		if err := r.DisplayChangefeed(); err != nil {
			return err
		}
	}
	return nil
}
