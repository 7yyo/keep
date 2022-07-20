package cdc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/list"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"keep/promp"
	net "keep/util/net"
	"keep/util/printer"
	"strings"
	"sync"
)

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

	body, err := r.Etcd.Get(context.TODO(), "/tidb/cdc/capture/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	r.captures = make([]capture, 0, len(body.Kvs))
	wg := &sync.WaitGroup{}
	wd := make(chan bool)
	defer close(wd)
	errors := make(chan error)
	defer close(errors)
	cm := make(map[string]capture)
	for _, kv := range body.Kvs {
		wg.Add(1)
		go r.hookCapture(kv, cm, wg, errors)
	}

	go func() {
		wg.Wait()
		wd <- true
	}()
	select {
	case <-wd:
		break
	case err := <-errors:
		return nil, err
	}

	if len(cm) == 0 {
		return nil, fmt.Errorf("no capture in this cluster")
	}
	return cm, nil
}

func (r *Runner) captureList() ([]string, error) {
	cs, err := r.captureInfo()
	if err != nil {
		return nil, err
	}
	cl := make([]string, 0, len(cs))
	for _, c := range cs {
		cl = append(cl, fmt.Sprintf("%s(%s)", c.Address, c.Id))
	}
	return cl, nil
}

func (r *Runner) hookCapture(kv *mvccpb.KeyValue, cMap map[string]capture, wg *sync.WaitGroup, errors chan error) {
	defer wg.Done()
	var c capture
	var cst captureStatus
	err := json.Unmarshal(kv.Value, &c)
	if err != nil {
		errors <- err
	}

	body, err := net.GetHttp(fmt.Sprintf("http://%s/api/v1/status", c.Address))
	if err != nil {
		errors <- err
	}
	err = json.Unmarshal(body, &cst)
	if err != nil {
		errors <- err
	}
	c.IsOwner = cst.IsOwner
	c.GitHash = cst.GitHash
	c.Pid = cst.Pid

	r.Lock()
	cMap[c.Address] = c
	r.captures = append(r.captures, c)
	r.Unlock()
}

func (r *Runner) displayCapture() error {
	cs, err := r.captureInfo()
	if err != nil {
		return err
	}
	captureOption := make([]string, 0, len(cs))
	captureOption = append(captureOption, printer.Return())
	for _, c := range cs {
		if c.IsOwner {
			captureOption = append(captureOption, fmt.Sprintf("%s (owner)", c.Address))
		} else {
			captureOption = append(captureOption, c.Address)
		}
	}
	p := promp.Select(captureOption, "cdc capture list", 20)
	_, c, err := p.Run()
	if err != nil {
		return err
	}
	if c == printer.Return() {
		if err := r.Run(); err != nil {
			return err
		}
	}
	c = strings.TrimSpace(strings.Split(c, "(")[0])
	l := list.NewWriter()
	l.SetStyle(list.StyleBulletCircle)
	l.AppendItems([]interface{}{
		fmt.Sprintf("id:      %s", cs[c].Id),
		fmt.Sprintf("owner:   %v", cs[c].IsOwner),
		fmt.Sprintf("version: %s", cs[c].Version),
		fmt.Sprintf("pid:     %d", cs[c].Pid)})
	fmt.Println(l.Render())
	return r.displayCapture()
	return nil
}
