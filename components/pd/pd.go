package components

import (
	"encoding/json"
	"fmt"
	net "keep/util/net"
)

type PlacementDriver struct {
	Header struct {
		ClusterID int64 `json:"cluster_id"`
	} `json:"header"`

	Members []struct {
		Name          string   `json:"name"`
		MemberID      int64    `json:"member_id"`
		PeerUrls      []string `json:"peer_urls"`
		ClientUrls    []string `json:"client_urls"`
		DeployPath    string   `json:"deploy_path"`
		BinaryVersion string   `json:"binary_version"`
		GitHash       string   `json:"git_hash"`
	} `json:"members"`

	Leader struct {
		Name          string   `json:"name"`
		MemberID      int64    `json:"member_id"`
		PeerUrls      []string `json:"peer_urls"`
		ClientUrls    []string `json:"client_urls"`
		DeployPath    string   `json:"deploy_path"`
		BinaryVersion string   `json:"binary_version"`
		GitHash       string   `json:"git_hash"`
	} `json:"leader"`

	EtcdLeader struct {
		Name          string   `json:"name"`
		MemberID      int64    `json:"member_id"`
		PeerUrls      []string `json:"peer_urls"`
		ClientUrls    []string `json:"client_urls"`
		DeployPath    string   `json:"deploy_path"`
		BinaryVersion string   `json:"binary_version"`
		GitHash       string   `json:"git_hash"`
	} `json:"etcd_leader"`
}

func NewPlacementDriver(pd string) *PlacementDriver {
	ps, err := net.GetHttp(fmt.Sprintf("http://%s/pd/api/v1/members", pd))
	if err != nil {
		panic(err)
	}
	var pdGroup PlacementDriver
	err = json.Unmarshal(ps, &pdGroup)
	return &pdGroup
}
