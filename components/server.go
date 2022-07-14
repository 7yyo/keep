package components

import (
	"github.com/manifoldco/promptui"
	"go.etcd.io/etcd/clientv3"
)

type Server struct {
	Etcd *clientv3.Client
	*PlacementDriver
}

func (s *Server) Cdc() error {

	prompt := promptui.Select{
		Label: " cdc ",
		Items: []string{
			"capture",
			"changefeed",
			"processor",
		},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return err
	}

	switch result {
	case "capture":
		captures, err := captureInfo(s.Etcd)
		if err != nil {
			return err
		}
		if err := displayCapture(captures); err != nil {
			return err
		}
	case "changefeed":
	default:

	}
	return nil
}
