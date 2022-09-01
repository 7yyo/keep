package cdc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/7yyo/sunflower/prompt"
	"github.com/jedib0t/go-pretty/v6/list"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	net "keep/util/net"
	"strings"
	"sync"
)

type capture struct {
	Id      string `json:"id"`
	Address string `json:"address"`
	captureStatus
}

type captureStatus struct {
	Version string `json:"version"`
	GitHash string `json:"git_hash"`
	Id      string `json:"id"`
	Pid     int    `json:"pid"`
	IsOwner bool   `json:"is_owner"`
}

func (r *Runner) getCaptureInfo() (map[string]capture, error) {
	body, err := r.Etcd.Get(context.TODO(), "/tidb/cdc/capture/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	wd := make(chan bool)
	defer close(wd)
	eh := make(chan error)
	defer close(eh)
	wg := &sync.WaitGroup{}
	cm := make(map[string]capture)
	for _, kv := range body.Kvs {
		wg.Add(1)
		go r.hookCapture(kv, cm, wg, eh)
	}
	go func() {
		wg.Wait()
		wd <- true
	}()
	select {
	case <-wd:
		break
	case err := <-eh:
		return nil, err
	}
	if len(cm) == 0 {
		return nil, fmt.Errorf("no capture in this cluster")
	}
	return cm, nil
}

func captureList(cm map[string]capture) []string {
	var captureList []string
	for _, c := range cm {
		captureList = append(captureList, fmt.Sprintf("%s(%s)", c.Address, c.Id))
	}
	return captureList
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

func (r *Runner) displayCaptureList() error {
	cmap, err := r.getCaptureInfo()
	if err != nil {
		return err
	}
	captureOption := make([]interface{}, 0, len(cmap))
	for _, c := range cmap {
		if c.IsOwner {
			captureOption = append(captureOption, fmt.Sprintf("%s (owner)", c.Address))
		} else {
			captureOption = append(captureOption, c.Address)
		}
	}
	se := prompt.Select{
		Title:  "capture:",
		Option: captureOption,
	}
	_, c, err := se.Run()
	if err != nil {
		if prompt.IsBackSpace(err) {
			if err := r.Run(); err != nil {
				return err
			}
		}
		return err
	}
	return r.displayCapture(cmap, c)
}

func (r *Runner) displayCapture(cm map[string]capture, c interface{}) error {
	c = strings.TrimSpace(strings.Split(c.(string), "(")[0])
	l := list.NewWriter()
	l.SetStyle(list.StyleBulletCircle)
	l.AppendItems([]interface{}{
		fmt.Sprintf("id:      %s", cm[c.(string)].Id),
		fmt.Sprintf("owner:   %v", cm[c.(string)].IsOwner),
		fmt.Sprintf("version: %s", cm[c.(string)].Version),
		fmt.Sprintf("pid:     %d", cm[c.(string)].Pid)})
	fmt.Println(l.Render())
	if prompt.Back() {
		return r.displayCaptureList()
	}
	return nil
}
