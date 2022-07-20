package cdc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/manifoldco/promptui"
	"github.com/tikv/client-go/v2/oracle"
	"golang.org/x/sync/errgroup"
	"keep/components/tidb"
	"keep/promp"
	"keep/util/color"
	net "keep/util/net"
	"keep/util/o"
	"keep/util/printer"
	"strconv"
	"strings"
	"sync"
	"time"
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

type changefeedConfig struct {
	UpstreamID int    `json:"upstream-id"`
	SinkURI    string `json:"sink-uri"`
	Opts       struct {
	} `json:"opts"`
	CreateTime   time.Time `json:"create-time"`
	StartTs      int64     `json:"start-ts"`
	TargetTs     int       `json:"target-ts"`
	AdminJobType int       `json:"admin-job-type"`
	SortEngine   string    `json:"sort-engine"`
	SortDir      string    `json:"sort-dir"`
	Config       struct {
		CaseSensitive    bool `json:"case-sensitive"`
		EnableOldValue   bool `json:"enable-old-value"`
		ForceReplicate   bool `json:"force-replicate"`
		CheckGcSafePoint bool `json:"check-gc-safe-point"`
		Filter           struct {
			Rules            []string    `json:"rules"`
			IgnoreTxnStartTs interface{} `json:"ignore-txn-start-ts"`
		} `json:"filter"`
		Mounter struct {
			WorkerNum int `json:"worker-num"`
		} `json:"mounter"`
		Sink struct {
			Dispatchers     interface{} `json:"dispatchers"`
			Protocol        string      `json:"protocol"`
			ColumnSelectors interface{} `json:"column-selectors"`
			SchemaRegistry  string      `json:"tpcc-registry"`
		} `json:"sink"`
		CyclicReplication struct {
			Enable           bool        `json:"enable"`
			ReplicaID        int         `json:"replica-id"`
			FilterReplicaIds interface{} `json:"filter-replica-ids"`
			IDBuckets        int         `json:"id-buckets"`
			SyncDdl          bool        `json:"sync-ddl"`
		} `json:"cyclic-replication"`
		Consistent struct {
			Level         string `json:"level"`
			MaxLogSize    int    `json:"max-log-size"`
			FlushInterval int    `json:"flush-interval"`
			Storage       string `json:"storage"`
		} `json:"consistent"`
	} `json:"config"`
	State             string      `json:"state"`
	Error             interface{} `json:"error"`
	SyncPointEnabled  bool        `json:"sync-point-enabled"`
	SyncPointInterval int64       `json:"sync-point-interval"`
	CreatorVersion    string      `json:"creator-version"`
}

func (r *Runner) changefeedInfo() (changefeed, error) {
	body, err := net.GetHttp(fmt.Sprintf("http://%s/api/v1/changefeeds", r.captures[0].Address))
	if err != nil {
		return nil, err
	}
	var cf changefeed
	if err = json.Unmarshal(body, &cf); err != nil {
		return nil, err
	}
	if len(cf) == 0 {
		fmt.Println(color.Red("no changefeed found, return..."))
		return nil, r.Run()
	}
	return cf, nil
}

func (r *Runner) changefeedInfoDetails() (*changefeedDetail, error) {
	body, err := net.GetHttp(fmt.Sprintf("http://%s/api/v1/changefeeds/%s", r.captures[0].Address, r.changefeedId))
	if err != nil {
		return nil, err
	}
	var cd changefeedDetail
	err = json.Unmarshal(body, &cd)
	if err != nil {
		return nil, err
	}
	return &cd, nil
}

func (r *Runner) displayChangefeedList() error {

	_, err := r.captureInfo()
	if err != nil {
		return err
	}
	cfs, err := r.changefeedInfo()
	if err != nil {
		return err
	}

	changefeedOption := make([]string, 0)
	changefeedOption = append(changefeedOption, printer.Return())
	for _, cf := range cfs {
		changefeedOption = append(changefeedOption, cf.Id)
	}
	p := promp.Select(changefeedOption, "cdc changefeed list", 20)
	_, c, err := p.Run()
	if err != nil {
		return err
	}
	r.changefeedId = c

	if c == printer.Return() {
		if err := r.Run(); err != nil {
			return err
		}
	}
	return r.displayChangefeedMenu()
}

var changefeedMenu = []string{
	printer.Return(),
	"info",
	"config",
	"pause",
	"resume",
	"transfer",
	"rebalance",
	color.Red("remove???"),
}

func (r *Runner) displayChangefeedMenu() error {

	p := promp.Select(changefeedMenu, "cdc changefeed menu", 20)
	i, _, err := p.Run()
	if err != nil {
		return err
	}

	switch i {
	case 0:
		err = r.displayChangefeedList()
	case 1:
		err = r.displayChangefeedDetails()
	case 2:
		err = r.displayChangefeedConfig()
	case 3:
		err = r.doChangefeed("pause", "stopped")
	case 4:
		err = r.doChangefeed("resume", "normal")
	case 5:
		err = r.transferTbl()
	case 6:
		err = r.curlChangefeed("rebalance")
	case 7:
		err = r.curlChangefeed("remove")
	}
	return err
}

func (r *Runner) displayChangefeedDetails() error {

	cd, err := r.changefeedInfoDetails()
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
	l.SetStyle(list.StyleConnectedLight)
	l.AppendItems([]interface{}{color.Green(fmt.Sprintf("[%s]", cd.ID))})
	l.Indent()
	l.AppendItems([]interface{}{
		fmt.Sprintf("sink_url:         %s", cd.SinkURI),
		fmt.Sprintf("create_time:      %s", cd.CreateTime),
		fmt.Sprintf("start_ts:         %d\n                  (%v)", cd.StartTs, oracle.GetTimeFromTS(uint64(cd.StartTs))),
		fmt.Sprintf("target_ts:        %d", cd.TargetTs),
		fmt.Sprintf("checkpoint_tso:   %d", cd.CheckpointTso),
		fmt.Sprintf("checkpoint_time:  %s", cd.CheckpointTime),
		fmt.Sprintf("sort_engine:      %s", cd.SortEngine),
		fmt.Sprintf("state:            %v", printer.Status(cd.State)),
		fmt.Sprintf("error:            %v", printer.IsNil(cd.Error)),
		"processor:"})
	l.Indent()

	cp := make(map[capture]processor)
	wg := &sync.WaitGroup{}
	wd := make(chan bool)
	errors := make(chan error)
	defer close(errors)
	for _, c := range r.captures {
		wg.Add(1)
		go r.hookProcessor(c, cp, wg, errors)
	}
	go func() {
		wg.Wait()
		close(wd)
	}()
	select {
	case <-wd:
		break
	case err := <-errors:
		if err != nil {
			return err
		}
	}

	pl, err := tidb.ListPartition(r.Pd.Leader.ClientUrls[0], o.TSOracle(), r.Etcd)
	if err != nil {
		return err
	}

	for c, pro := range cp {
		l.AppendItems([]interface{}{
			fmt.Sprintf("%s\n%s", color.Green(c.Address), color.Green(c.Id)),
		})
		l.Indent()
		l.AppendItems([]interface{}{
			fmt.Sprintf("checkpoint_ts: %d\n               (%s)", pro.CheckpointTs, oracle.GetTimeFromTS(uint64(pro.CheckpointTs))),
			fmt.Sprintf("resolved_ts:   %d\n               (%s)", pro.ResolvedTs, oracle.GetTimeFromTS(uint64(pro.ResolvedTs))),
			fmt.Sprintf("lag:           %.2fs", oracle.GetTimeFromTS(uint64(pro.ResolvedTs)).Sub(oracle.GetTimeFromTS(uint64(pro.CheckpointTs))).Seconds()),
			"table list:",
		})
		l.Indent()
		tblStr := ""
		for _, tbl := range pro.TableIds {
			for _, p := range pl {
				if p.ID == tbl {
					tblStr += fmt.Sprintf("[`%s`.`%s`.`%s`(%s)], ", p.DBName, p.Table.Name.O, p.Name, strconv.FormatInt(p.ID, 10))
					break
				} else if p.Table.ID == tbl {
					tblStr += fmt.Sprintf("[%s]`%s`.`%s`, ", strconv.FormatInt(p.Table.ID, 10), p.DBName, p.Table.Name.O)
					break
				}
			}
		}
		in := make([]interface{}, 0)
		in = append(in, tblStr)
		l.AppendItems(in)
		l.UnIndent()
		l.UnIndent()
	}
	fmt.Println(l.Render())
	return r.displayChangefeedMenu()
}

func (r *Runner) doChangefeed(c string, t string) error {
	p := promptui.Prompt{
		Label:     "return",
		IsConfirm: true,
	}
	result, _ := p.Run()
	if result != "y" {
		return r.displayChangefeedMenu()
	} else {
		req := fmt.Sprintf("http://%s/api/v1/changefeeds/%s/%s", r.captures[0].Address, r.changefeedId, c)
		if _, err := net.PostHttp(req, ""); err != nil {
			return err
		}
		group := new(errgroup.Group)
		group.Go(func() error {
			err := r.changefeedState(t)
			if err != nil {
				return err
			}
			return nil
		})
		if err := group.Wait(); err != nil {
			return err
		}
		return r.displayChangefeedMenu()
	}
}

func (r *Runner) changefeedState(state string) error {
	ticker := time.NewTicker(time.Second * 1)
	retry := 0
	for range ticker.C {
		cd, err := r.changefeedInfoDetails()
		if err != nil {
			return err
		}
		if cd.State == state {
			ticker.Stop()
			fmt.Printf(color.Green("changefeed state changed to [%s] complete\n"), state)
			break
		} else {
			retry++
			if retry > 10 {
				ticker.Stop()
				return fmt.Errorf("%s changefeed failed after retry more than 10 times, please check why\n", state)
			}
		}
	}
	return nil
}

func (r *Runner) displayChangefeedConfig() error {
	d, err := r.Etcd.Get(context.TODO(), fmt.Sprintf("/tidb/cdc/changefeed/info/%s", r.changefeedId))
	if err != nil {
		return err
	}
	var cc changefeedConfig
	if len(d.Kvs) != 0 {
		if err := json.Unmarshal(d.Kvs[0].Value, &cc); err != nil {
			return err
		}
	}
	l := list.NewWriter()
	l.SetStyle(list.StyleConnectedLight)
	l.AppendItems([]interface{}{"config"})
	l.Indent()
	l.AppendItems([]interface{}{
		fmt.Sprintf("caseSensitive: %v", cc.Config.CaseSensitive),
		fmt.Sprintf("enableOldValue: %v", cc.Config.EnableOldValue),
		fmt.Sprintf("forceReplicate: %v", cc.Config.ForceReplicate),
		fmt.Sprintf("checkGcSafePoint: %v", cc.Config.CheckGcSafePoint),
		"filter: ",
	})
	l.Indent()
	l.AppendItems([]interface{}{
		fmt.Sprintf("rules: %v", cc.Config.Filter.Rules),
		fmt.Sprintf("ignoreTxnStartTs: %v", cc.Config.Filter.IgnoreTxnStartTs),
	})
	l.UnIndent()
	l.AppendItems([]interface{}{
		"mounter: ",
	})
	l.Indent()
	l.AppendItems([]interface{}{
		fmt.Sprintf("workerNum: %d", cc.Config.Mounter.WorkerNum),
	})
	l.UnIndent()
	l.AppendItems([]interface{}{
		"sink",
	})
	l.Indent()
	l.AppendItems([]interface{}{
		fmt.Sprintf("dispatchers: %v", cc.Config.Sink.Dispatchers),
		fmt.Sprintf("protocol: %v", cc.Config.Sink.Protocol),
		fmt.Sprintf("columnSelectors: %v", cc.Config.Sink.ColumnSelectors),
		fmt.Sprintf("schemaRegistry: %v", cc.Config.Sink.SchemaRegistry),
	})
	l.UnIndent()
	l.AppendItems([]interface{}{
		"cyclicReplication",
	})
	l.Indent()
	l.AppendItems([]interface{}{
		fmt.Sprintf("enable: %v", cc.Config.CyclicReplication.Enable),
		fmt.Sprintf("replicaID: %d", cc.Config.CyclicReplication.ReplicaID),
		fmt.Sprintf("filterReplicaIds: %v", cc.Config.CyclicReplication.FilterReplicaIds),
		fmt.Sprintf("IDBuckets: %d", cc.Config.CyclicReplication.IDBuckets),
		fmt.Sprintf("syncDdl: %v", cc.Config.CyclicReplication.SyncDdl),
	})
	l.UnIndent()
	l.AppendItems([]interface{}{
		"consistent",
	})
	l.Indent()
	l.AppendItems([]interface{}{
		fmt.Sprintf("level: %v", cc.Config.Consistent.Level),
		fmt.Sprintf("maxLogSize: %v", cc.Config.Consistent.MaxLogSize),
		fmt.Sprintf("flushInterval: %v", cc.Config.Consistent.FlushInterval),
		fmt.Sprintf("storage: %v", cc.Config.Consistent.Storage),
	})
	fmt.Println(l.Render())
	return r.displayChangefeedMenu()
}

func (r *Runner) transferTbl() error {

	pl, err := tidb.ListPartition(r.Pd.Leader.ClientUrls[0], o.CheckPointTs(), r.Etcd)
	if err != nil {
		return err
	}
	partitions := make([]string, 0, len(pl))
	partitions = append(partitions, printer.Return())
	for _, p := range pl {
		if p.Name != "" {
			partitions = append(partitions, fmt.Sprintf("[%d]`%s`.`%s`.`%s`", p.ID, p.DBName, p.Table.Name, p.Name))
		} else {
			partitions = append(partitions, fmt.Sprintf("[%d]`%s`.`%s`", p.Table.ID, p.DBName, p.Table.Name))
		}
	}

	searcher := func(input string, index int) bool {
		pName := partitions[index]
		name := strings.Replace(strings.ToLower(pName), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
	p := promp.Select(partitions, "choose table to transfer", 20)
	p.Searcher = searcher
	i, tbl, err := p.Run()
	if i == 0 {
		return r.displayChangefeedMenu()
	}

	tblId, err := strconv.Atoi(strings.Split(tbl, "]")[0][1:])
	if err != nil {
		return err
	}

	cl, err := r.captureList()
	if err != nil {
		return err
	}

	ps := promp.Select(cl, "choose capture to transfer", 20)
	_, cpId, err := ps.Run()
	if err != nil {
		return err
	}

	if i == 0 {
		return r.displayChangefeedMenu()
	}
	postData := map[string]interface{}{
		"capture_id": strings.Split(cpId, "(")[1][:len(strings.Split(cpId, "(")[1])-1],
		"table_id":   tblId,
	}

	fmt.Println()
	pp := promptui.Prompt{
		Label:     "transfer?",
		IsConfirm: true,
	}
	result, _ := pp.Run()
	if result != "y" {
		return r.displayChangefeedMenu()
	}

	req := fmt.Sprintf("http://%s/api/v1/changefeeds/%s/tables/move_table", r.captures[0].Address, r.changefeedId)
	if err := net.Curl(req, postData); err != nil {
		return err
	} else {
		fmt.Println(fmt.Sprintf(color.Green("%s =>>> %s success"), tbl, cpId))
	}
	return r.displayChangefeedMenu()
}

func (r *Runner) curlChangefeed(c string) error {
	pp := promptui.Prompt{
		Label:     c + "?",
		IsConfirm: true,
	}
	result, _ := pp.Run()
	if result != "y" {
		return r.displayChangefeedMenu()
	}

	var req string
	switch c {
	case "remove":
		req = fmt.Sprintf("http://%s/api/v1/changefeeds/%s", r.captures[0].Address, r.changefeedId)
		if err := net.DeleteCurl(req); err != nil {
			return err
		} else {
			fmt.Println(fmt.Sprintf(color.Green("%sd %s"), c, r.changefeedId))
		}
		return r.displayChangefeedList()
	case "rebalance":
		req = fmt.Sprintf("http://%s/api/v1/changefeeds/%s/tables/rebalance_table", r.captures[0].Address, r.changefeedId)
		if err := net.Curl(req, nil); err != nil {
			return err
		} else {
			fmt.Println(fmt.Sprintf(color.Green("%s `%s` success"), c, r.changefeedId))
			return r.displayChangefeedMenu()
		}
	default:
		return nil
	}
}
