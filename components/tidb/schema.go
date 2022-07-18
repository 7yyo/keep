package tidb

import (
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/pingcap/tidb/meta"
	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tiflow/cdc/kv"
	"github.com/pingcap/tiflow/pkg/config"
	"go.etcd.io/etcd/clientv3"
	net "keep/util/net"
	"keep/util/o"
	"strconv"
	"strings"
	"sync"
)

type DB struct {
	Id     int64
	Name   string
	Tables []*model.TableInfo
}

type Partition struct {
	ID     int64
	Name   string
	Table  *model.TableInfo
	DBName string
}

func snapshotMeta(endpoint string, checkPointTS uint64) (*meta.Meta, error) {
	conf := config.GetGlobalServerConfig()
	storage, err := kv.CreateTiStore(endpoint, conf.Security)
	if err != nil {
		return nil, err
	}
	return kv.GetSnapshotMeta(storage, checkPointTS)
}

func ListDatabase(endpoint string, checkPointTs uint64) (*[]DB, error) {
	snapshot, err := snapshotMeta(endpoint, checkPointTs)
	if err != nil {
		return nil, err
	}
	ld, err := snapshot.ListDatabases()
	if err != nil {
		return nil, err
	}
	dbl := make([]DB, 0)
	for _, d := range ld {
		if d.Name.String() == "mysql" {
			continue
		}
		ts, err := snapshot.ListTables(d.ID)
		if err != nil {
			return nil, err
		}
		db := DB{
			Id:     d.ID,
			Name:   d.Name.String(),
			Tables: ts,
		}
		dbl = append(dbl, db)
	}
	return &dbl, nil
}

func ListPartition(endpoint string, checkPointTs uint64, etcd *clientv3.Client) ([]Partition, error) {
	ss, err := snapshotMeta(endpoint, checkPointTs)
	if err != nil {
		return nil, err
	}
	ld, err := ss.ListDatabases()
	if err != nil {
		return nil, err
	}
	tidbCluster := NewTiDBCluster(etcd)
	partitions := make([]Partition, 0)
	errors := make(chan error)
	defer close(errors)
	wg := &sync.WaitGroup{}
	for _, d := range ld {
		wg.Add(1)
		go hookPartition(&partitions, tidbCluster, d, ss, wg, errors)
	}
	wg.Wait()
	return partitions, err
}

func hookPartition(partitions *[]Partition, tidbCluster []*TiDB, d *model.DBInfo, meta *meta.Meta, wg *sync.WaitGroup, errors chan error) {
	defer wg.Done()
	if d.Name.String() == "mysql" {
		return
	}
	tbls, err := meta.ListTables(d.ID)
	if err != nil {
		errors <- err
	}
	for _, tbl := range tbls {
		tblInfo, err := tblInfo(d.Name.String(), tbl.Name.String(), tidbCluster[0].Host, tidbCluster[0].StatusPort)
		if err != nil {
			errors <- err
		}
		if tblInfo.Partition != nil {
			for _, definition := range tblInfo.Partition.Definitions {
				partition := Partition{
					ID:     definition.ID,
					Name:   definition.Name.O,
					Table:  tbl,
					DBName: d.Name.O,
				}
				*partitions = append(*partitions, partition)
			}
		} else {
			partition := Partition{
				Table:  tbl,
				DBName: d.Name.O,
			}
			*partitions = append(*partitions, partition)
		}
	}
}

func (r *Runner) displayTiDBSchema() error {
	dbs, err := ListDatabase(r.Pd.Leader.ClientUrls[0], o.CheckPointTs())
	if err != nil {
		return err
	}
	tblOption := make([]string, 0)
	tblOption = append(tblOption, "return?")
	for _, db := range *dbs {
		for _, tbl := range db.Tables {
			tblOption = append(tblOption, fmt.Sprintf("`%s`.`%s`", db.Name, tbl.Name.String()))
		}
	}
	p := promptui.Select{
		Label: "tidb table list",
		Items: tblOption,
		Size:  20,
	}
	_, c, err := p.Run()

	if c == "return?" {
		if err := r.Run(); err != nil {
			return err
		}
	}
	return nil
}

func tblInfo(db string, tbl string, host string, statusPort int) (*model.TableInfo, error) {
	req := fmt.Sprintf("http://%s:%s/schema/%s/%s", strings.Split(host, ":")[0], strconv.Itoa(statusPort), db, tbl)
	body, err := net.GetHttp(req)
	if err != nil {
		return nil, err
	}
	tblInfo := &model.TableInfo{}
	err = json.Unmarshal(body, tblInfo)
	return tblInfo, err
}
