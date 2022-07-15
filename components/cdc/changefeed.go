package cdc

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/manifoldco/promptui"
	"github.com/tikv/client-go/v2/oracle"
	"keep/components/tidb"
	"keep/util/color"
	net "keep/util/net"
	o "keep/util/oracle"
	"keep/util/printer"
	"keep/util/sys"
	"sync"
)

type changefeed []struct {
	Id             string      `json:"id"`
	State          string      `json:"state"`
	CheckpointTso  int64       `json:"checkpoint_tso"`
	CheckpointTime string      `json:"checkpoint_time"`
	Error          interface{} `json:"error"`
}

type changefeedDetail struct {
	ID             string      `json:"id"`
	SinkURI        string      `json:"sink_uri"`
	CreateTime     string      `json:"create_time"`
	StartTs        int64       `json:"start_ts"`
	TargetTs       int         `json:"target_ts"`
	CheckpointTso  uint64      `json:"checkpoint_tso"`
	CheckpointTime string      `json:"checkpoint_time"`
	SortEngine     string      `json:"sort_engine"`
	State          string      `json:"state"`
	Error          interface{} `json:"error"`
	ErrorHistory   interface{} `json:"error_history"`
	CreatorVersion string      `json:"creator_version"`
	TaskStatus     []struct {
		CaptureID       string  `json:"capture_id"`
		TableIds        []int64 `json:"table_ids"`
		TableOperations struct {
		} `json:"table_operations"`
	} `json:"task_status"`
}

func changefeedInfo(h string) (changefeed, error) {
	r, err := net.GetHttp(fmt.Sprintf("http://%s/api/v1/changefeeds", h))
	if err != nil {
		return nil, err
	}
	var cf changefeed
	if err = json.Unmarshal(r, &cf); err != nil {
		return nil, err
	}
	if len(cf) == 0 {
		return nil, fmt.Errorf("no changefeed found")
	}
	return cf, nil
}

func (r *Runner) displayChangefeed() error {

	_, err := r.captureInfo()
	if err != nil {
		return err
	}
	cfs, err := changefeedInfo(r.captures[0].Address)
	if err != nil {
		return err
	}

	changefeedOption := make([]string, 0)
	for _, cf := range cfs {
		changefeedOption = append(changefeedOption, cf.Id)
	}
	changefeedOption = append(changefeedOption, "return?")

	p := promptui.Select{
		Label: "cdc changefeed list",
		Items: changefeedOption,
	}
	_, c, err := p.Run()
	if err != nil {
		return err
	}

	if c == "return?" {
		if err := r.Run(); err != nil {
			return err
		}
	}

	r.changefeedId = c
	return r.displayChangefeedMenu()
}

func (r *Runner) displayChangefeedMenu() error {

	p := promptui.Select{
		Label: r.changefeedId,
		Items: []string{
			"details?",
			"return?",
		},
	}
	_, c, err := p.Run()
	if err != nil {
		return err
	}

	switch c {
	case "details?":
		if err := r.displayChangefeedDetails(); err != nil {
			return err
		}
	case "return?":
		if err := r.displayChangefeed(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) displayChangefeedDetails() error {
	body, err := net.GetHttp(fmt.Sprintf("http://%s/api/v1/changefeeds/%s", r.captures[0].Address, r.changefeedId))
	if err != nil {
		return err
	}
	var cd changefeedDetail
	err = json.Unmarshal(body, &cd)
	if err != nil {
		return err
	}

	partitions, err := tidb.ListPartition(r.Pd.Leader.ClientUrls[0], o.CheckPointTs(), r.Etcd)
	if err != nil {
		return err
	}
	pm := make(map[int8]string)
	for _, partition := range partitions {
		if partition.Name == "" {
			pm[int8(partition.Table.ID)] = fmt.Sprintf("`%s`.`%s`", partition.DBName, partition.Table.Name.O)
		} else {
			pm[int8(partition.ID)] = fmt.Sprintf("`%s`.`%s` [%s]", partition.DBName, partition.Table.Name.O, partition.Name)
		}
	}

	l := list.NewWriter()
	l.SetStyle(list.StyleBulletCircle)
	l.AppendItems([]interface{}{color.Green(fmt.Sprintf("[%s]", cd.ID))})
	l.Indent()
	l.AppendItems([]interface{}{
		fmt.Sprintf("sink_url:         %s", cd.SinkURI),
		fmt.Sprintf("create_time:      %s", cd.CreateTime),
		fmt.Sprintf("start_ts:         %d", cd.StartTs),
		fmt.Sprintf("target_ts:        %d", cd.TargetTs),
		fmt.Sprintf("checkpoint_tso:   %d", cd.CheckpointTso),
		fmt.Sprintf("checkpoint_time:  %s", cd.CheckpointTime),
		fmt.Sprintf("sort_engine:      %s", cd.SortEngine),
		fmt.Sprintf("state:            %v", printer.Status(cd.State)),
		fmt.Sprintf("error:            %v", printer.IsNil(cd.Error)),
		fmt.Sprintf("error_history:    %v", printer.IsNil(cd.ErrorHistory)),
		"",
		"processor: "})
	l.Indent()

	cm := make(map[capture]processor)
	wg := &sync.WaitGroup{}
	errors := make(chan error)
	defer close(errors)
	for _, c := range r.captures {
		wg.Add(1)
		go r.hookProcessor(c, cm, wg, errors)
	}
	wg.Wait()

	pl, err := tidb.ListPartition(r.Pd.Leader.ClientUrls[0], o.TSOracle(), r.Etcd)
	if err != nil {
		return err
	}

	for c, pro := range cm {
		l.AppendItems([]interface{}{
			fmt.Sprintf("%s\n%s", color.Green(c.Address), color.Green(c.Id)),
		})
		l.Indent()
		l.AppendItems([]interface{}{
			fmt.Sprintf("checkpoint_ts: %d (%s)", pro.CheckpointTs, oracle.GetTimeFromTS(uint64(pro.CheckpointTs))),
			fmt.Sprintf("resolved_ts:   %d (%s)", pro.ResolvedTs, oracle.GetTimeFromTS(uint64(pro.ResolvedTs))),
			fmt.Sprintf("lag:           %.2fs", oracle.GetTimeFromTS(uint64(pro.ResolvedTs)).Sub(oracle.GetTimeFromTS(uint64(pro.CheckpointTs))).Seconds()),
			"tables:",
		})
		l.Indent()
		for _, tbl := range pro.TableIds {
			for _, p := range pl {
				if p.ID == int64(tbl) {
					l.AppendItems([]interface{}{
						fmt.Sprintf("`%s`.`%s` [%s]", p.DBName, p.Table.Name.O, p.Name),
					})
					break
				}
				if p.Table.ID == int64(tbl) {
					l.AppendItems([]interface{}{
						fmt.Sprintf("`%s`.`%s`", p.DBName, p.Table.Name.O),
					})
					break
				}
			}
		}
		l.UnIndent()
		l.UnIndent()
	}
	fmt.Println(l.Render())

	prompt := promptui.Prompt{
		Label:     "return",
		IsConfirm: true,
	}
	result, _ := prompt.Run()
	if result == "y" {
		return r.displayChangefeedMenu()
	} else {
		sys.Exit()
	}
	return nil
}
