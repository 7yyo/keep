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
	"keep/util/color"
)

func main() {

	defer fmt.Println("\ngoodbye!")
	endpoints := args()

	etcd := components.NewEtcd(endpoints)
	pdGroup := pd.NewPlacementDriver(endpoints[0])
	tidbCluster := tidb.NewTiDBCluster(etcd)

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
	_, m, _ := prompt.Run()

	log.SetLevel(zapcore.PanicLevel)
	if err := s.Run(m); err != nil {
		fmt.Println(color.Red(fmt.Sprintf("sorry, an error occurred: \n%s", err.Error())))
	}

}

func args() []string {
	var endpoint string
	flag.StringVar(&endpoint, "pd", "", "")
	flag.Parse()

	if endpoint == "" {
		endpoint = "172.16.5.133:2379"
	}
	endpoints := make([]string, 0)
	endpoints = append(endpoints, endpoint)
	return endpoints
}
