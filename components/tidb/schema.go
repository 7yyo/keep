package tidb

import (
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tiflow/cdc/kv"
	"github.com/pingcap/tiflow/pkg/config"
	"keep/util"
	"keep/util/time"
	"strconv"
	"strings"
)

type Schema struct {
	DBs []DB
}

type DB struct {
	Id     int
	Name   string
	Tables []*model.TableInfo
}

func SchemaInfo(endpoint string, checkPointTSO uint64) (*Schema, error) {
	conf := config.GetGlobalServerConfig()
	storage, err := kv.CreateTiStore(endpoint, conf.Security)
	if err != nil {
		return nil, err
	}
	snapshot, err := kv.GetSnapshotMeta(storage, checkPointTSO)
	if err != nil {
		return nil, err
	}
	dbs, err := snapshot.ListDatabases()
	if err != nil {
		return nil, err
	}
	var databases []DB
	var database DB
	for _, db := range dbs {
		if db.Name.String() == "mysql" {
			continue
		}
		database.Name = db.Name.String()
		database.Id = int(db.ID)
		tls, err := snapshot.ListTables(db.ID)
		if err != nil {
			return nil, err
		}
		database.Tables = tls
		databases = append(databases, database)
	}
	sm := Schema{
		DBs: databases,
	}
	return &sm, nil
}

func (r *Runner) displayTiDBSchema() (string, error) {
	sm, err := SchemaInfo(r.Pd.Leader.ClientUrls[0], time.CheckPointTSO())
	if err != nil {
		return "", err
	}
	tblOption := make([]string, 0)
	for _, db := range sm.DBs {
		for _, tbl := range db.Tables {
			tblOption = append(tblOption, fmt.Sprintf("%s.%s", db.Name, tbl.Name.String()))
		}
	}
	p := promptui.Select{
		Label: "tidb table list",
		Items: tblOption,
	}
	_, tbl, err := p.Run()
	if err != nil {
		return "", err
	}
	return tbl, nil
}

func TableInfo(id string, host string, statusPort int) (*model.TableInfo, error) {
	db := strings.Split(id, ".")[0]
	tbl := strings.Split(id, ".")[1]
	url := "http://" + strings.Split(host, ":")[0] + ":" + strconv.Itoa(statusPort) + "/" + "schema" + "/" + db + "/" + tbl
	body, err := util.GetHttp(url)
	if err != nil {
		return nil, err
	}
	tblInfo := &model.TableInfo{}
	err = json.Unmarshal(body, tblInfo)
	return tblInfo, err
}
