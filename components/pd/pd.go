package components

import (
	"encoding/json"
	"fmt"
	"github.com/7yyo/sunflower/prompt"
	"go.etcd.io/etcd/clientv3"
	"io"
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

type config struct {
	Name     string
	Default  string
	Current  string
	Optional string
	Unit     string
}

type Runner struct {
	Pd   *PlacementDriver
	Etcd *clientv3.Client
}

var pdOptions = []interface{}{
	" config",
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

func (r *Runner) Run() error {
	p := prompt.Select{
		Option: pdOptions,
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

func (r *Runner) configs() ([]interface{}, error) {

	cfg, err := r.getConfig(false)
	if err != nil {
		return nil, err
	}
	dCfg, err := r.getConfig(true)
	if err != nil {
		return nil, err
	}

	configs := make([]interface{}, 0)
	configs = append(configs, config{Name: "cluster-version", Current: cfg.ClusterVersion, Default: dCfg.ClusterVersion})
	configs = append(configs, config{Name: "log.level", Current: cfg.Log.Level, Default: dCfg.Log.Level, Optional: "debug, info, warn, error, fatal"})
	configs = append(configs, config{Name: "schedule.max-merge-region-size", Current: strconv.Itoa(cfg.Schedule.MaxMergeRegionSize), Default: strconv.Itoa(dCfg.Schedule.MaxMergeRegionSize), Unit: "MiB"})
	configs = append(configs, config{Name: "schedule.max-merge-region-keys", Current: strconv.Itoa(dCfg.Schedule.MaxMergeRegionKeys), Default: strconv.Itoa(cfg.Schedule.MaxMergeRegionKeys)})
	configs = append(configs, config{Name: "schedule.patrol-region-interval", Current: cfg.Schedule.PatrolRegionInterval, Default: dCfg.Schedule.PatrolRegionInterval})
	configs = append(configs, config{Name: "schedule.split-merge-interval", Current: cfg.Schedule.SplitMergeInterval, Default: dCfg.Schedule.SplitMergeInterval})
	configs = append(configs, config{Name: "schedule.max-snapshot-count", Current: strconv.Itoa(cfg.Schedule.MaxSnapshotCount), Default: strconv.Itoa(dCfg.Schedule.MaxSnapshotCount)})
	configs = append(configs, config{Name: "schedule.max-pending-peer-count", Current: strconv.Itoa(cfg.Schedule.MaxPendingPeerCount), Default: strconv.Itoa(dCfg.Schedule.MaxPendingPeerCount)})
	configs = append(configs, config{Name: "schedule.max-store-down-time", Current: cfg.Schedule.MaxStoreDownTime, Default: dCfg.Schedule.MaxStoreDownTime})
	configs = append(configs, config{Name: "schedule.leader-schedule-policy", Current: cfg.Schedule.LeaderSchedulePolicy, Default: dCfg.Schedule.LeaderSchedulePolicy})
	configs = append(configs, config{Name: "schedule.leader-schedule-limit", Current: strconv.Itoa(cfg.Schedule.LeaderScheduleLimit), Default: strconv.Itoa(dCfg.Schedule.LeaderScheduleLimit)})
	configs = append(configs, config{Name: "schedule.region-schedule-limit", Current: strconv.Itoa(cfg.Schedule.RegionScheduleLimit), Default: strconv.Itoa(dCfg.Schedule.RegionScheduleLimit)})
	configs = append(configs, config{Name: "schedule.replica-schedule-limit", Current: strconv.Itoa(cfg.Schedule.ReplicaScheduleLimit), Default: strconv.Itoa(dCfg.Schedule.ReplicaScheduleLimit)})
	configs = append(configs, config{Name: "schedule.merge-schedule-limit", Current: strconv.Itoa(cfg.Schedule.MergeScheduleLimit), Default: strconv.Itoa(dCfg.Schedule.MergeScheduleLimit)})
	configs = append(configs, config{Name: "schedule.hot-region-schedule-limit", Current: strconv.Itoa(cfg.Schedule.HotRegionScheduleLimit), Default: strconv.Itoa(dCfg.Schedule.HotRegionScheduleLimit)})
	configs = append(configs, config{Name: "schedule.hot-region-cache-hits-threshold", Current: strconv.Itoa(cfg.Schedule.HotRegionCacheHitsThreshold), Default: strconv.Itoa(dCfg.Schedule.HotRegionCacheHitsThreshold)})
	configs = append(configs, config{Name: "schedule.high-space-ratio", Current: strconv.FormatFloat(cfg.Schedule.HighSpaceRatio, 'g', 12, 64), Default: strconv.FormatFloat(dCfg.Schedule.HighSpaceRatio, 'g', 12, 64)})
	configs = append(configs, config{Name: "schedule.low-space-ratio", Current: strconv.FormatFloat(cfg.Schedule.LowSpaceRatio, 'g', 12, 64), Default: strconv.FormatFloat(dCfg.Schedule.LowSpaceRatio, 'g', 12, 64)})
	configs = append(configs, config{Name: "schedule.tolerant-size-ratio", Current: strconv.Itoa(cfg.Schedule.TolerantSizeRatio), Default: strconv.Itoa(dCfg.Schedule.TolerantSizeRatio)})
	configs = append(configs, config{Name: "schedule.enable-remove-down-replica", Current: cfg.Schedule.EnableRemoveDownReplica, Default: dCfg.Schedule.EnableRemoveDownReplica})
	configs = append(configs, config{Name: "schedule.enable-replace-offline-replica", Current: cfg.Schedule.EnableReplaceOfflineReplica, Default: dCfg.Schedule.EnableReplaceOfflineReplica})
	configs = append(configs, config{Name: "schedule.enable-make-up-replica", Current: cfg.Schedule.EnableMakeUpReplica, Default: dCfg.Schedule.EnableMakeUpReplica})
	configs = append(configs, config{Name: "schedule.enable-remove-extra-replica", Current: cfg.Schedule.EnableRemoveExtraReplica, Default: dCfg.Schedule.EnableRemoveExtraReplica})
	configs = append(configs, config{Name: "schedule.enable-location-replacement", Current: cfg.Schedule.EnableLocationReplacement, Default: dCfg.Schedule.EnableLocationReplacement})
	configs = append(configs, config{Name: "schedule.enable-cross-table-merge", Current: cfg.Schedule.EnableCrossTableMerge, Default: dCfg.Schedule.EnableCrossTableMerge})
	configs = append(configs, config{Name: "schedule.enable-one-way-merge", Current: cfg.Schedule.EnableOneWayMerge, Default: dCfg.Schedule.EnableOneWayMerge})
	configs = append(configs, config{Name: "replication.max-replicas", Current: strconv.Itoa(cfg.Replication.MaxReplicas), Default: strconv.Itoa(dCfg.Replication.MaxReplicas)})
	configs = append(configs, config{Name: "replication.location-labels", Current: cfg.Replication.LocationLabels, Default: dCfg.Replication.LocationLabels})
	configs = append(configs, config{Name: "replication.enable-placement-rules", Current: cfg.Replication.EnablePlacementRules, Default: dCfg.Replication.EnablePlacementRules})
	configs = append(configs, config{Name: "replication.strictly-match-label", Current: cfg.Replication.StrictlyMatchLabel, Default: dCfg.Replication.StrictlyMatchLabel})
	configs = append(configs, config{Name: "pd-server.use-region-storage", Current: cfg.PdServer.UseRegionStorage, Default: dCfg.PdServer.UseRegionStorage})
	configs = append(configs, config{Name: "pd-server.max-gap-reset-ts", Current: cfg.PdServer.MaxGapResetTs, Default: dCfg.PdServer.MaxGapResetTs})
	configs = append(configs, config{Name: "pd-server.key-type", Current: cfg.PdServer.KeyType, Default: dCfg.PdServer.KeyType})
	configs = append(configs, config{Name: "pd-server.metric-storage", Current: cfg.PdServer.MetricStorage, Default: dCfg.PdServer.MetricStorage})
	configs = append(configs, config{Name: "pd-server.dashboard-address", Current: cfg.PdServer.DashboardAddress, Default: dCfg.PdServer.DashboardAddress})
	configs = append(configs, config{Name: "replication-mode.replication-mode", Current: cfg.ReplicationMode.ReplicationMode, Default: dCfg.ReplicationMode.ReplicationMode, Optional: "majority, dr-auto-sync"})
	return configs, nil
}

func (r *Runner) getConfig(d bool) (*Config, error) {
	req := ""
	if !d {
		req = fmt.Sprintf("%s/pd/api/v1/config", r.Pd.Leader.ClientUrls[0])
	} else {
		req = fmt.Sprintf("%s/pd/api/v1/config/default", r.Pd.Leader.ClientUrls[0])
	}

	body, err := net.GetHttp(req)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err = json.Unmarshal(body, &cfg); err != nil {
		return nil, err
	}
	return &cfg, err
}

func (r *Runner) displayConfigs() error {

	configs, err := r.configs()
	if err != nil {
		return err
	}

	s := prompt.Select{
		Option: configs,
		Title:  "Select config to display",
		Desc: &prompt.Desc{
			T: "Name",
			D: []string{
				"Current",
				"Default",
				"Optional",
				"Unit",
			},
		},
		Cap: 15,
	}
	i, _, err := s.Run()
	if prompt.IsBackSpace(err) {
		return r.Run()
	}

	cfg := configs[i].(config)
	v, err := prompt.Assign(fmt.Sprintf("set %s =", cfg.Name))
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
		jsonBody = fmt.Sprintf("{\"%s\":%d}", cfg.Name, cp)
	} else {
		cp, err := strconv.ParseFloat(v, 64)
		if err != nil {
			jsonBody = fmt.Sprintf("{\"%s\":\"%s\"}", cfg.Name, v)
		} else {
			jsonBody = fmt.Sprintf("{\"%s\":%.1f}", cfg.Name, cp)
		}
	}

	var resp *http.Response
	resp, _ = net.PostHttp(req, jsonBody)
	if resp.StatusCode != 200 {
		message, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		printer.PrintError(fmt.Sprintf("update failed, error: %s", string(message)))
	} else {
		fmt.Println("updated!")
	}
	return r.displayConfigs()
}
