package main

import (
	"flag"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/pingcap/log"
	"go.uber.org/zap/zapcore"
	"keep/components"
	pd "keep/components/pd"
	"keep/components/tidb"
)

func main() {

	defer fmt.Println("\ngoodbye!")

	var err error
	endpoints := args()

	etcd := components.NewEtcd(endpoints)
	pdGroup := pd.NewPlacementDriver(endpoints[0])
	tidbCluster := tidb.NewTiDBCluster(*etcd)

	s := components.Server{
		Etcd:            etcd,
		PlacementDriver: pdGroup,
		TiDBCluster:     tidbCluster,
	}

	prompt := promptui.Select{
		Label: "what can i do for you",
		Items: []string{
			"cdc",
			"tidb",
		},
	}
	_, m, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	log.SetLevel(zapcore.PanicLevel)
	if err := s.Run(m); err != nil {
		panic(err)
	}

}

func args() []string {
	var endpoint string
	flag.StringVar(&endpoint, "pd", "", "")
	flag.Parse()

	if endpoint == "" {
		endpoint = "172.16.5.133:2379"
	}
	var endpoints []string
	endpoints = append(endpoints, endpoint)
	return endpoints
}
