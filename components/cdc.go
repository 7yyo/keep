package components

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"go.etcd.io/etcd/clientv3"
	"keep/util"
	"strings"
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

func captureInfo(etcd *clientv3.Client) ([]capture, error) {
	r, err := etcd.Get(context.TODO(), "/tidb/cdc/capture/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var c capture
	var cs []capture
	var cStatus captureStatus
	for _, kv := range r.Kvs {
		err = json.Unmarshal(kv.Value, &c)
		if err != nil {
			return nil, err
		}
		body, err := util.GetHttp(fmt.Sprintf("http://%s/api/v1/status", c.Address))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &cStatus)
		if err != nil {
			return nil, err
		}
		c.IsOwner = cStatus.IsOwner
		c.GitHash = cStatus.GitHash
		c.Pid = cStatus.Pid
		cs = append(cs, c)
	}
	return cs, nil
}

func displayCapture(captures []capture) error {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F449 {{ .Address | cyan }} ",
		Inactive: "{{ .Address | cyan }} ",
		Selected: "\U0001F449 {{ .Address | red | cyan }}",
		Details: `
				{{ "Id:" | faint }}	{{ .Id }}
				{{ "Address:" | faint }}	{{ .Address }}
				{{ "Version:" | faint }}	{{ .Version }}
				{{ "Is_Owner:" | faint }}	{{ .IsOwner }}
				{{ "Pid:" | faint }}	{{ .Pid }}
				{{ "Git_Hash:" | faint }}	{{ .GitHash }}`,
	}
	searcher := func(input string, index int) bool {
		c := captures[index]
		name := strings.Replace(strings.ToLower(c.Address), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "choose capture",
		Items:     captures,
		Templates: templates,
		Size:      5,
		Searcher:  searcher,
	}

	_, _, err := prompt.Run()
	if err != nil {
		return err
	}
	return nil
}
