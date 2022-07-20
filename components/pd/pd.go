package components

import (
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"io"
	"keep/util/color"
	net "keep/util/net"
	"keep/util/printer"
	"net/http"
	"strconv"
	"strings"
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

type Config struct {
	ClientUrls          string `json:"client-urls"`
	PeerUrls            string `json:"peer-urls"`
	AdvertiseClientUrls string `json:"advertise-client-urls"`
	AdvertisePeerUrls   string `json:"advertise-peer-urls"`
	Name                string `json:"name"`
	DataDir             string `json:"data-dir"`
	ForceNewCluster     bool   `json:"force-new-cluster"`
	EnableGrpcGateway   bool   `json:"enable-grpc-gateway"`
	InitialCluster      string `json:"initial-cluster"`
	InitialClusterState string `json:"initial-cluster-state"`
	InitialClusterToken string `json:"initial-cluster-token"`
	Join                string `json:"join"`
	Lease               int    `json:"lease"`
	Log                 struct {
		Level            string `json:"level"`
		Format           string `json:"format"`
		DisableTimestamp bool   `json:"disable-timestamp"`
		File             struct {
			Filename   string `json:"filename"`
			MaxSize    int    `json:"max-size"`
			MaxDays    int    `json:"max-days"`
			MaxBackups int    `json:"max-backups"`
		} `json:"file"`
		Development         bool        `json:"development"`
		DisableCaller       bool        `json:"disable-caller"`
		DisableStacktrace   bool        `json:"disable-stacktrace"`
		DisableErrorVerbose bool        `json:"disable-error-verbose"`
		Sampling            interface{} `json:"sampling"`
	} `json:"log"`
	TsoSaveInterval           string `json:"tso-save-interval"`
	TsoUpdatePhysicalInterval string `json:"tso-update-physical-interval"`
	EnableLocalTso            bool   `json:"enable-local-tso"`
	Metric                    struct {
		Job      string `json:"job"`
		Address  string `json:"address"`
		Interval string `json:"interval"`
	} `json:"metric"`
	Schedule struct {
		MaxSnapshotCount            int    `json:"max-snapshot-count"`
		MaxPendingPeerCount         int    `json:"max-pending-peer-count"`
		MaxMergeRegionSize          int    `json:"max-merge-region-size"`
		MaxMergeRegionKeys          int    `json:"max-merge-region-keys"`
		SplitMergeInterval          string `json:"split-merge-interval"`
		EnableOneWayMerge           string `json:"enable-one-way-merge"`
		EnableCrossTableMerge       string `json:"enable-cross-table-merge"`
		PatrolRegionInterval        string `json:"patrol-region-interval"`
		MaxStoreDownTime            string `json:"max-store-down-time"`
		MaxStorePreparingTime       string `json:"max-store-preparing-time"`
		LeaderScheduleLimit         int    `json:"leader-schedule-limit"`
		LeaderSchedulePolicy        string `json:"leader-schedule-policy"`
		RegionScheduleLimit         int    `json:"region-schedule-limit"`
		ReplicaScheduleLimit        int    `json:"replica-schedule-limit"`
		MergeScheduleLimit          int    `json:"merge-schedule-limit"`
		HotRegionScheduleLimit      int    `json:"hot-region-schedule-limit"`
		HotRegionCacheHitsThreshold int    `json:"hot-region-cache-hits-threshold"`
		StoreLimit                  struct {
			Num1 struct {
				AddPeer    int `json:"add-peer"`
				RemovePeer int `json:"remove-peer"`
			} `json:"1"`
		} `json:"store-limit"`
		TolerantSizeRatio           int     `json:"tolerant-size-ratio"`
		LowSpaceRatio               float64 `json:"low-space-ratio"`
		HighSpaceRatio              float64 `json:"high-space-ratio"`
		RegionScoreFormulaVersion   string  `json:"region-score-formula-version"`
		SchedulerMaxWaitingOperator int     `json:"scheduler-max-waiting-operator"`
		EnableRemoveDownReplica     string  `json:"enable-remove-down-replica"`
		EnableReplaceOfflineReplica string  `json:"enable-replace-offline-replica"`
		EnableMakeUpReplica         string  `json:"enable-make-up-replica"`
		EnableRemoveExtraReplica    string  `json:"enable-remove-extra-replica"`
		EnableLocationReplacement   string  `json:"enable-location-replacement"`
		EnableDebugMetrics          string  `json:"enable-debug-metrics"`
		EnableJointConsensus        string  `json:"enable-joint-consensus"`
		SchedulersV2                []struct {
			Type        string      `json:"type"`
			Args        interface{} `json:"args"`
			Disable     bool        `json:"disable"`
			ArgsPayload string      `json:"args-payload"`
		} `json:"schedulers-v2"`
		SchedulersPayload struct {
			BalanceHotRegionScheduler interface{} `json:"balance-hot-region-scheduler"`
			BalanceLeaderScheduler    struct {
				Batch  int `json:"batch"`
				Ranges []struct {
					EndKey   string `json:"end-key"`
					StartKey string `json:"start-key"`
				} `json:"ranges"`
			} `json:"balance-leader-scheduler"`
			BalanceRegionScheduler struct {
				Name   string `json:"name"`
				Ranges []struct {
					EndKey   string `json:"end-key"`
					StartKey string `json:"start-key"`
				} `json:"ranges"`
			} `json:"balance-region-scheduler"`
			SplitBucketScheduler interface{} `json:"split-bucket-scheduler"`
		} `json:"schedulers-payload"`
		StoreLimitMode          string `json:"store-limit-mode"`
		HotRegionsWriteInterval string `json:"hot-regions-write-interval"`
		HotRegionsReservedDays  int    `json:"hot-regions-reserved-days"`
	} `json:"schedule"`
	Replication struct {
		MaxReplicas               int    `json:"max-replicas"`
		LocationLabels            string `json:"location-labels"`
		StrictlyMatchLabel        string `json:"strictly-match-label"`
		EnablePlacementRules      string `json:"enable-placement-rules"`
		EnablePlacementRulesCache string `json:"enable-placement-rules-cache"`
		IsolationLevel            string `json:"isolation-level"`
	} `json:"replication"`
	PdServer struct {
		UseRegionStorage                 string `json:"use-region-storage"`
		MaxGapResetTs                    string `json:"max-gap-reset-ts"`
		KeyType                          string `json:"key-type"`
		RuntimeServices                  string `json:"runtime-services"`
		MetricStorage                    string `json:"metric-storage"`
		DashboardAddress                 string `json:"dashboard-address"`
		TraceRegionFlow                  string `json:"trace-region-flow"`
		FlowRoundByDigit                 int    `json:"flow-round-by-digit"`
		MinResolvedTsPersistenceInterval string `json:"min-resolved-ts-persistence-interval"`
	} `json:"pd-server"`
	ClusterVersion string `json:"cluster-version"`
	Labels         struct {
	} `json:"labels"`
	QuotaBackendBytes         string `json:"quota-backend-bytes"`
	AutoCompactionMode        string `json:"auto-compaction-mode"`
	AutoCompactionRetentionV2 string `json:"auto-compaction-retention-v2"`
	TickInterval              string `json:"TickInterval"`
	ElectionInterval          string `json:"ElectionInterval"`
	PreVote                   bool   `json:"PreVote"`
	MaxRequestBytes           int    `json:"max-request-bytes"`
	Security                  struct {
		CacertPath    string      `json:"cacert-path"`
		CertPath      string      `json:"cert-path"`
		KeyPath       string      `json:"key-path"`
		CertAllowedCn interface{} `json:"cert-allowed-cn"`
		SSLCABytes    interface{} `json:"SSLCABytes"`
		SSLCertBytes  interface{} `json:"SSLCertBytes"`
		SSLKEYBytes   interface{} `json:"SSLKEYBytes"`
		RedactInfoLog bool        `json:"redact-info-log"`
		Encryption    struct {
			DataEncryptionMethod  string `json:"data-encryption-method"`
			DataKeyRotationPeriod string `json:"data-key-rotation-period"`
			MasterKey             struct {
				Type     string `json:"type"`
				KeyID    string `json:"key-id"`
				Region   string `json:"region"`
				Endpoint string `json:"endpoint"`
				Path     string `json:"path"`
			} `json:"master-key"`
		} `json:"encryption"`
	} `json:"security"`
	LabelProperty struct {
	} `json:"label-property"`
	WarningMsgs                 interface{} `json:"WarningMsgs"`
	DisableStrictReconfigCheck  bool        `json:"DisableStrictReconfigCheck"`
	HeartbeatStreamBindInterval string      `json:"HeartbeatStreamBindInterval"`
	LeaderPriorityCheckInterval string      `json:"LeaderPriorityCheckInterval"`
	Dashboard                   struct {
		TidbCacertPath     string `json:"tidb-cacert-path"`
		TidbCertPath       string `json:"tidb-cert-path"`
		TidbKeyPath        string `json:"tidb-key-path"`
		PublicPathPrefix   string `json:"public-path-prefix"`
		InternalProxy      bool   `json:"internal-proxy"`
		EnableTelemetry    bool   `json:"enable-telemetry"`
		EnableExperimental bool   `json:"enable-experimental"`
	} `json:"dashboard"`
	ReplicationMode struct {
		ReplicationMode string `json:"replication-mode"`
		DrAutoSync      struct {
			LabelKey         string `json:"label-key"`
			Primary          string `json:"primary"`
			Dr               string `json:"dr"`
			PrimaryReplicas  int    `json:"primary-replicas"`
			DrReplicas       int    `json:"dr-replicas"`
			WaitStoreTimeout string `json:"wait-store-timeout"`
			WaitSyncTimeout  string `json:"wait-sync-timeout"`
			WaitAsyncTimeout string `json:"wait-async-timeout"`
			PauseRegionSplit string `json:"pause-region-split"`
		} `json:"dr-auto-sync"`
	} `json:"replication-mode"`
}

type Runner struct {
	Pd *PlacementDriver
}

func (r *Runner) Run() error {
	p := promptui.Select{
		Label: "placement driver",
		Items: []string{
			"config",
		},
	}
	i, _, err := p.Run()
	if err != nil {
		return err
	}
	switch i {
	case 0:
		return r.displayConfigs()
	}
	return nil
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

func (r *Runner) configs() (map[string]string, []string, error) {
	req := fmt.Sprintf("%s/pd/api/v1/config", r.Pd.Leader.ClientUrls[0])
	body, err := net.GetHttp(req)
	if err != nil {
		return nil, nil, err
	}
	var cfg Config
	if err = json.Unmarshal(body, &cfg); err != nil {
		return nil, nil, err
	}
	cm := make(map[string]string)
	cm["cluster-version"] = cfg.ClusterVersion
	cm["log.level"] = cfg.Log.Level
	cm["schedule.max-merge-region-size"] = strconv.Itoa(cfg.Schedule.MaxMergeRegionSize)
	cm["schedule.max-merge-region-keys"] = strconv.Itoa(cfg.Schedule.MaxMergeRegionKeys)
	cm["schedule.patrol-region-interval"] = cfg.Schedule.PatrolRegionInterval
	cm["schedule.split-merge-interval"] = cfg.Schedule.SplitMergeInterval
	cm["schedule.max-snapshot-count"] = strconv.Itoa(cfg.Schedule.MaxSnapshotCount)
	cm["schedule.max-pending-peer-count"] = strconv.Itoa(cfg.Schedule.MaxPendingPeerCount)
	cm["schedule.max-store-down-time"] = cfg.Schedule.MaxStoreDownTime
	cm["schedule.leader-schedule-policy"] = cfg.Schedule.LeaderSchedulePolicy
	cm["schedule.leader-schedule-limit"] = strconv.Itoa(cfg.Schedule.LeaderScheduleLimit)
	cm["schedule.region-schedule-limit"] = strconv.Itoa(cfg.Schedule.RegionScheduleLimit)
	cm["schedule.replica-schedule-limit"] = strconv.Itoa(cfg.Schedule.ReplicaScheduleLimit)
	cm["schedule.merge-schedule-limit"] = strconv.Itoa(cfg.Schedule.MergeScheduleLimit)
	cm["schedule.hot-region-schedule-limit"] = strconv.Itoa(cfg.Schedule.HotRegionScheduleLimit)
	cm["schedule.hot-region-cache-hits-threshold"] = strconv.Itoa(cfg.Schedule.HotRegionCacheHitsThreshold)
	cm["schedule.high-space-ratio"] = strconv.FormatFloat(cfg.Schedule.HighSpaceRatio, 'g', 12, 64)
	cm["schedule.low-space-ratio"] = strconv.FormatFloat(cfg.Schedule.LowSpaceRatio, 'g', 12, 64)
	cm["schedule.tolerant-size-ratio"] = strconv.Itoa(cfg.Schedule.TolerantSizeRatio)
	cm["schedule.enable-remove-down-replica"] = cfg.Schedule.EnableRemoveDownReplica
	cm["schedule.enable-replace-offline-replica"] = cfg.Schedule.EnableReplaceOfflineReplica
	cm["schedule.enable-make-up-replica"] = cfg.Schedule.EnableMakeUpReplica
	cm["schedule.enable-remove-extra-replica"] = cfg.Schedule.EnableRemoveExtraReplica
	cm["schedule.enable-location-replacement"] = cfg.Schedule.EnableLocationReplacement
	cm["schedule.enable-cross-table-merge"] = cfg.Schedule.EnableCrossTableMerge
	cm["schedule.enable-one-way-merge"] = cfg.Schedule.EnableOneWayMerge
	cm["replication.max-replicas"] = strconv.Itoa(cfg.Replication.MaxReplicas)
	cm["replication.location-labels"] = cfg.Replication.LocationLabels
	cm["replication.enable-placement-rules"] = cfg.Replication.EnablePlacementRules
	cm["replication.strictly-match-label"] = cfg.Replication.StrictlyMatchLabel
	cm["pd-server.use-region-storage"] = cfg.PdServer.UseRegionStorage
	cm["pd-server.max-gap-reset-ts"] = cfg.PdServer.MaxGapResetTs
	cm["pd-server.key-type"] = cfg.PdServer.KeyType
	cm["pd-server.metric-storage"] = cfg.PdServer.MetricStorage
	cm["pd-server.dashboard-address"] = cfg.PdServer.DashboardAddress
	cm["replication-mode.replication-mode"] = cfg.ReplicationMode.ReplicationMode

	keys := make([]string, 0, len(cm))
	keys = append(keys, printer.Return())
	for k, _ := range cm {
		keys = append(keys, k)
	}
	return cm, keys, nil
}

func (r *Runner) displayConfigs() error {
	m, c, err := r.configs()
	if err != nil {
		return err
	}
	searcher := func(input string, index int) bool {
		pName := c[index]
		name := strings.Replace(strings.ToLower(pName), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
	p := promptui.Select{
		Label:    "configs",
		Items:    c,
		Searcher: searcher,
		Size:     20,
	}
	_, result, err := p.Run()
	if err != nil {
		return err
	}
	if result == printer.Return() {
		return r.Run()
	}
	pp := promptui.Prompt{
		Label: fmt.Sprintf("current `%s`, type `%s` return", color.Green(m[result]), "\\q"),
	}
	v, err := pp.Run()
	if err != nil {
		return err
	}
	if v == "\\q" || strings.TrimSpace(v) == "" {
		return r.displayConfigs()
	}
	req := fmt.Sprintf("%s/pd/api/v1/config", r.Pd.Leader.ClientUrls[0])

	cp, err := strconv.Atoi(v)
	var jsonBody string
	if err == nil {
		jsonBody = fmt.Sprintf("{\"%s\":%d}", result, cp)
	} else {
		cp, err := strconv.ParseFloat(v, 64)
		if err != nil {
			jsonBody = fmt.Sprintf("{\"%s\":\"%s\"}", result, v)
		} else {
			jsonBody = fmt.Sprintf("{\"%s\":%.1f}", result, cp)
		}
	}

	var resp *http.Response
	resp, _ = net.PostHttp(req, jsonBody)
	if resp.StatusCode != 200 {
		message, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		printer.PrintError(fmt.Sprintf("update failed, error:\n%s", string(message)))
	} else {
		fmt.Println(color.Green("update success!"))
	}
	return r.displayConfigs()
}
