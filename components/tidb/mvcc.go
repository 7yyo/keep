package tidb

import (
	"encoding/base64"
	"fmt"
	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/tablecodec"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tidb/util/rowcodec"
	"time"
)

type mvcc struct {
	Key      string `json:"key"`
	RegionID int    `json:"region_id"`
	Value    struct {
		Info struct {
			Writes []struct {
				StartTs    int    `json:"start_ts"`
				CommitTs   int    `json:"commit_ts"`
				ShortValue string `json:"short_value"`
			} `json:"writes"`
		} `json:"info"`
	} `json:"value"`
}

func decodeMVCC(tblInfo *model.TableInfo, snapshot string) (map[string]string, error) {
	bs, err := base64.StdEncoding.DecodeString(snapshot)
	if err != nil {
		return nil, err
	}
	colMap := make(map[int64]*types.FieldType, 3)
	for _, col := range tblInfo.Columns {
		colMap[col.ID] = &col.FieldType
	}
	r, err := decodeRow(bs, colMap, time.UTC)
	m := make(map[string]string)
	for _, col := range tblInfo.Columns {
		if v, ok := r[col.ID]; ok {
			if v.IsNull() {
				m[col.Name.L] = "NULL"
				continue
			}
			ss, err := v.ToString()
			if err != nil {
				m[col.Name.L] = "error: " + err.Error() + fmt.Sprintf("datum: %#v", v)
				continue
			}
			m[col.Name.L] = ss
		}
	}
	return m, nil
}

func decodeRow(b []byte, cols map[int64]*types.FieldType, loc *time.Location) (map[int64]types.Datum, error) {
	if !rowcodec.IsNewFormat(b) {
		return tablecodec.DecodeRowWithMap(b, cols, loc, nil)
	}
	return tablecodec.DecodeRowWithMapNew(b, cols, loc, nil)
}
