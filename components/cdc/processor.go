package cdc

import (
	"encoding/json"
	"fmt"
	net "keep/util/net"
	"sync"
)

type processor struct {
	CheckpointTs int64       `json:"checkpoint_ts"`
	ResolvedTs   int64       `json:"resolved_ts"`
	TableIds     []int64     `json:"table_ids"`
	Count        int         `json:"count"`
	Error        interface{} `json:"error"`
}

func (r *Runner) hookProcessor(c capture, cm map[capture]processor, wg *sync.WaitGroup, errors chan error) {
	defer wg.Done()
	body, err := net.GetHttp(fmt.Sprintf("http://%s/api/v1/processors/%s/%s", c.Address, r.changefeedId, c.Id))
	if err != nil {
		errors <- err
	}
	var pro processor
	if err := json.Unmarshal(body, &pro); err != nil {
		errors <- err
	}
	r.Lock()
	cm[c] = pro
	r.Unlock()
}
