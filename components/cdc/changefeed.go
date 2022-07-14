package cdc

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/manifoldco/promptui"
	"keep/components/tidb"
	"keep/util"
	"keep/util/printer"
	"keep/util/sys"
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
	r, err := util.GetHttp(fmt.Sprintf("http://%s/api/v1/changefeeds", h))
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

func (r *Runner) DisplayChangefeed() error {

	_, err := r.captureInfo()
	if err != nil {
		return err
	}
	cfs, err := changefeedInfo(r.captures[0].Address)
	if err != nil {
		return err
	}

	var changefeedOption []string
	for _, cf := range cfs {
		changefeedOption = append(changefeedOption, cf.Id)
	}
	changefeedOption = append(changefeedOption, "return?")

	p := promptui.Select{
		Label: "changefeed list",
		Items: changefeedOption,
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

	r.changefeedId = m
	if err := r.displayChangefeedDetail(); err != nil {
		return err
	}

	return nil
}

func (r *Runner) displayChangefeedDetail() error {
	body, err := util.GetHttp(fmt.Sprintf("http://%s/api/v1/changefeeds/%s", r.captures[0].Address, r.changefeedId))
	if err != nil {
		return err
	}
	var cd changefeedDetail
	err = json.Unmarshal(body, &cd)
	if err != nil {
		return err
	}

	sm, err := tidb.SchemaInfo(r.Pd.Leader.ClientUrls[0], cd.CheckpointTso)
	if err != nil {
		return err
	}
	sMap := make(map[int64]string)
	for _, s := range sm.DBs {
		for _, t := range s.Tables {
			sMap[t.ID] = fmt.Sprintf("`%s`.`%s`", s.Name, t.Name.String())
		}
	}

	l := list.NewWriter()
	l.SetStyle(list.StyleBulletCircle)
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
		"task_status: "})
	l.Indent()
	for _, ts := range cd.TaskStatus {
		for _, cp := range r.captures {
			if cp.Id == ts.CaptureID {
				l.AppendItems([]interface{}{fmt.Sprintf("%s", cp.Address)})
				break
			}
		}
		if len(ts.TableIds) > 0 {
			l.Indent()
			l.Indent()
			for _, t := range ts.TableIds {
				l.AppendItems([]interface{}{fmt.Sprintf("%v", sMap[t])})
			}
			l.UnIndent()
		}
	}
	fmt.Println(l.Render())

	prompt := promptui.Prompt{
		Label:     "return",
		IsConfirm: true,
	}
	result, _ := prompt.Run()
	if result == "y" {
		if err := r.DisplayChangefeed(); err != nil {
			return err
		}
	} else {
		sys.Exit()
	}

	return nil
}
