package cdc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/manifoldco/promptui"
	"go.etcd.io/etcd/clientv3"
	"keep/util"
	"keep/util/sys"
)

const captureUrl string = "/tidb/cdc/capture/"

type capture struct {
	Id      string `json:"id"`
	Address string `json:"address"`
	Version string `json:"version"`
	IsOwner bool
	Pid     int
	GitHash string
}

type captureStatus struct {
	Version string `json:"version"`
	GitHash string `json:"git_hash"`
	Id      string `json:"id"`
	Pid     int    `json:"pid"`
	IsOwner bool   `json:"is_owner"`
}

func (r *Runner) captureInfo() (map[string]capture, error) {
	body, err := r.Etcd.Get(context.TODO(), captureUrl, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var c capture
	cMap := make(map[string]capture)
	var cst captureStatus
	for _, kv := range body.Kvs {
		err = json.Unmarshal(kv.Value, &c)
		if err != nil {
			return nil, err
		}
		body, err := util.GetHttp(fmt.Sprintf("http://%s/api/v1/status", c.Address))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &cst)
		if err != nil {
			return nil, err
		}
		c.IsOwner = cst.IsOwner
		c.GitHash = cst.GitHash
		c.Pid = cst.Pid
		cMap[c.Address] = c
		r.captures = append(r.captures, c)
	}
	if len(cMap) == 0 {
		return nil, fmt.Errorf("no capture")
	}
	return cMap, nil
}

func (r *Runner) DisplayCapture() error {

	cMap, err := r.captureInfo()
	if err != nil {
		return err
	}

	var captureOption []string
	for _, c := range cMap {
		captureOption = append(captureOption, c.Address)
	}
	captureOption = append(captureOption, "return?")

	p := promptui.Select{
		Label: "capture list",
		Items: captureOption,
	}
	_, m, err := p.Run()
	if err != nil {
		return err
	}

	if m == "return?" {
		if err := r.Run(); err != nil {
			return err
		}
	}

	l := list.NewWriter()
	l.SetStyle(list.StyleBulletCircle)
	l.AppendItems([]interface{}{
		fmt.Sprintf("id: %s", cMap[m].Id),
		fmt.Sprintf("owner: %v", cMap[m].IsOwner),
		fmt.Sprintf("version: %s", cMap[m].Version),
		fmt.Sprintf("pid: %d", cMap[m].Pid)})
	fmt.Println(l.Render())

	prompt := promptui.Prompt{
		Label:     "return",
		IsConfirm: true,
	}
	result, _ := prompt.Run()
	if result == "y" {
		if err := r.DisplayCapture(); err != nil {
			return err
		}
	} else {
		sys.Exit()
	}
	return nil
}
